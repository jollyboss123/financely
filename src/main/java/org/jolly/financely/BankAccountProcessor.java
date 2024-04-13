package org.jolly.financely;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.function.Function;

/**
 * @author jolly
 */
@Component
public class BankAccountProcessor implements ItemProcessor<RawTransaction, Transaction> {
    private static final Logger log = LoggerFactory.getLogger(BankAccountProcessor.class);
    private DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("ddMMMyy");
    private int dateLen = 7;
    private Function<String, Boolean> isCreditTransfer;
    private final DefaultFieldExtractor transferAmountExtractor;
    private final DefaultFieldExtractor instalmentExtractor;

    public BankAccountProcessor(DefaultFieldExtractor transferAmountExtractor, DefaultFieldExtractor instalmentExtractor) {
        this.transferAmountExtractor = transferAmountExtractor;
        this.instalmentExtractor = instalmentExtractor;
    }

    public void setDateTimeFormatter(DateTimeFormatter dateTimeFormatter) {
        this.dateTimeFormatter = dateTimeFormatter;
    }

    public void setDateLen(int dateLen) {
        this.dateLen = dateLen;
    }

    public void setIsCreditTransfer(Function<String, Boolean> isCreditTransfer) {
        this.isCreditTransfer = isCreditTransfer;
    }

    @Override
    public Transaction process(@NonNull RawTransaction item) {
        final LocalDate date = extractDate(item);
        String fullDesc = item.getMergedLines(dateLen);

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
        if (Boolean.TRUE.equals(isCreditTransfer.apply(desc))) {
            credit = Long.parseLong(amountStr);
        } else {
            debit = Long.parseLong(amountStr);
        }

        return new Transaction.Builder(item.getFile(),1L, date, Bank.valueOf(MDC.get(MDCKey.BANK.name())), desc)
                .credit(credit)
                .debit(debit)
                .instalment(instalment)
                .isInstalment(isInstalment)
                .build();
    }

    private LocalDate extractDate(RawTransaction rawTransaction) {
        String dateString = rawTransaction.getLines().getFirst().substring(0,dateLen);
        return LocalDate.parse(dateString, dateTimeFormatter);
    }
}
