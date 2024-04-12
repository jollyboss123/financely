package org.jolly.financely;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

/**
 * @author jolly
 */
@Component
public class BankAccountProcessor implements ItemProcessor<RawTransaction, Transaction> {
    private static final Logger log = LoggerFactory.getLogger(BankAccountProcessor.class);
    private DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern("ddMMMyy");
    private int dateLen = 7;
    private final DefaultFieldExtractor transferAmountExtractor;

    public BankAccountProcessor(DefaultFieldExtractor transferAmountExtractor) {
        this.transferAmountExtractor = transferAmountExtractor;
    }

    public void setDateTimeFormatter(DateTimeFormatter dateTimeFormatter) {
        this.dateTimeFormatter = dateTimeFormatter;
    }

    public void setDateLen(int dateLen) {
        this.dateLen = dateLen;
    }

    @Override
    public Transaction process(@NonNull RawTransaction item) {
        final LocalDate date = extractDate(item);
        final String desc = item.getMergedLines(dateLen);
        final String amountStr = transferAmountExtractor.getField(desc)
                .replace(",", "")
                .replace(".", "");
        long credit = 0;
        long debit = 0;
        if (desc.contains("CR")) {
            credit = Long.parseLong(amountStr);
        } else {
            debit = Long.parseLong(amountStr);
        }
        return new Transaction.Builder(item.getFile(),1L, date, "UOB", desc)
                .credit(credit)
                .debit(debit)
                .build();
    }

    private LocalDate extractDate(RawTransaction rawTransaction) {
        String dateString = rawTransaction.getLines().getFirst().substring(0,dateLen);
        return LocalDate.parse(dateString, dateTimeFormatter);
    }
}