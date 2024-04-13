package org.jolly.financely.batch.processor;

import org.jolly.financely.batch.extractor.DefaultFieldExtractor;
import org.jolly.financely.constant.Bank;
import org.jolly.financely.constant.MDCKey;
import org.jolly.financely.model.Instalment;
import org.jolly.financely.model.RawTransaction;
import org.jolly.financely.model.Transaction;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.beans.factory.config.ConfigurableBeanFactory;
import org.springframework.context.annotation.Scope;
import org.springframework.data.util.Pair;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;

/**
 * @author jolly
 */
@Component(value = "BankAccountProcessor")
@Scope(value = ConfigurableBeanFactory.SCOPE_PROTOTYPE)
public class BankAccountProcessor implements ItemProcessor<RawTransaction, Transaction> {
    private static final Logger log = LoggerFactory.getLogger(BankAccountProcessor.class);
    private DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("ddMMMyy");
    // allow for optional min, max length for date strings
    // set max to null to disable optional
    private DateLength dateLengths;
    private String[] creditTransfer;
    private String[] itemsToSkip;
    private final DefaultFieldExtractor transferAmountExtractor;
    private final DefaultFieldExtractor instalmentExtractor;

    public BankAccountProcessor(DefaultFieldExtractor transferAmountExtractor, DefaultFieldExtractor instalmentExtractor) {
        this.transferAmountExtractor = transferAmountExtractor;
        this.instalmentExtractor = instalmentExtractor;
    }

    public void setDateTimeFormatter(DateTimeFormatter dateTimeFormatter) {
        this.dateTimeFormatter = dateTimeFormatter;
    }

    public void setDateLengths(DateLength dateLengths) {
        this.dateLengths = dateLengths;
    }

    public void setCreditTransfer(String[] creditTransfer) {
        this.creditTransfer = creditTransfer;
    }

    public void setItemsToSkip(String[] itemsToSkip) {
        this.itemsToSkip = itemsToSkip;
    }

    @Override
    public Transaction process(@NonNull RawTransaction item) {
        final DateInfo dateInfo = extractDate(item);
        String fullDesc = item.getMergedLines(dateInfo.length());

        if (shouldSkip(fullDesc)) {
            return null;
        }

        final String instalmentStr = instalmentExtractor.getField(fullDesc);
        Instalment instalment = null;
        boolean isInstalment = false;
        if (StringUtils.hasText(instalmentStr)) {
            fullDesc = fullDesc.replace(instalmentStr, "");
            String[] s = instalmentStr.split("/");
            instalment = new Instalment.Builder()
                    .number(Integer.valueOf(s[0]))
                    .total(Integer.valueOf(s[1]))
                    .build();
            isInstalment = true;
        }

        final String amountStr = transferAmountExtractor.getField(fullDesc)
                .replace(",", "")
                .replace(".", "");
        final String desc = fullDesc.replace(transferAmountExtractor.getField(fullDesc), "");
        long credit = 0;
        long debit = 0;
        if (isCreditTransfer(desc)) {
            credit = Long.parseLong(amountStr);
        } else {
            debit = Long.parseLong(amountStr);
        }

        return new Transaction.Builder(item.getFile(),1L, dateInfo.date, Bank.valueOf(MDC.get(MDCKey.BANK.name())), desc)
                .credit(credit)
                .debit(debit)
                .instalment(instalment)
                .isInstalment(isInstalment)
                .build();
    }

    private DateInfo extractDate(RawTransaction rawTransaction) {
        int i = dateLengths.min();
        String dateStr = null;

        if (dateLengths.max() == null) {
            dateStr = rawTransaction.getLines().getFirst().substring(0, i);
            final LocalDate localDate = LocalDate.parse(dateStr, dateTimeFormatter);
            return new DateInfo(i, localDate);
        } else {
            while (i <= dateLengths.max()) {
                dateStr = rawTransaction.getLines().getFirst().substring(0, i);
                try {
                    final LocalDate localDate = LocalDate.parse(dateStr, dateTimeFormatter);
                    return new DateInfo(i, localDate);
                } catch (DateTimeParseException ignored) {
                    // continue
                }
                i++;
            }
        }

        throw new DateTimeParseException("Date not found or invalid format", dateStr, i);
    }

    private boolean isCreditTransfer(String desc) {
        if (creditTransfer != null) {
            for (String ct : creditTransfer) {
                if (desc.matches(ct)) {
                    return true;
                }
            }
        }
        return false;
    }

    private boolean shouldSkip(String desc) {
        if (desc.trim().isEmpty()) {
            return true;
        }

        if (itemsToSkip != null) {
            for (String item : itemsToSkip) {
                if (desc.matches(item)) {
                    return true;
                }
            }
        }

        return false;
    }

    public record DateLength(int min, Integer max){}
    private record DateInfo(int length, LocalDate date){}
}
