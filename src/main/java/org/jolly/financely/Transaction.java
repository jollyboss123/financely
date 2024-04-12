package org.jolly.financely;

import java.text.SimpleDateFormat;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

/**
 * @author jolly
 */
public final class Transaction implements Comparable<Transaction> {
    private final long id;
    private final LocalDate date;
    private final long debit;
    private final long credit;
    private final String head;
    private final String subHead;
    private final String description;
    private final String file;
    private final boolean isSalary;
    private final String bank;
    private static final DateTimeFormatter DATE_FORMAT = DateTimeFormatter.ofPattern("yyyy-MM-dd");

    public static class Builder {
        // required params
        private final long id;
        private final LocalDate date;
        private final String bank;
        private final String description;
        private final String file;

        // optional params
        private long debit = 0;
        private long credit = 0;
        private String head;
        private String subHead;
        private boolean isSalary = false;

        public Builder(String file, long id, LocalDate date, String bank, String description) {
            this.file = file;
            this.id = id;
            this.date = date;
            this.bank = bank;
            this.description = description;
        }

        public Builder debit(long val) {
            debit = val;
            return this;
        }

        public Builder credit(long val) {
            credit = val;
            return this;
        }

        public Builder head(String val) {
            head = val;
            return this;
        }

        public Builder subHead(String val) {
            subHead = val;
            return this;
        }

        public Builder isSalary(boolean val) {
            isSalary = val;
            return this;
        }

        public Transaction build() {
            return new Transaction(this);
        }
    }

    private Transaction(Builder builder) {
        this.file = builder.file;
        this.id = builder.id;
        this.bank = builder.bank;
        this.description = builder.description;
        this.date = builder.date;
        this.debit = builder.debit;
        this.credit = builder.credit;
        this.head = builder.head;
        this.subHead = builder.subHead;
        this.isSalary = builder.isSalary;
    }

    public String getDateStr() {
        return DATE_FORMAT.format(date);
    }

    @Override
    public int hashCode() {
        int result = Long.hashCode(id);
        result = 31 * result + file.hashCode();
        result = 31 * result + (date != null ? date.hashCode() : 0);
        result = 31 * result + (Long.hashCode(debit));
        result = 31 * result + (Long.hashCode(credit));
        result = 31 * result + (head != null ? head.hashCode() : 0);
        result = 31 * result + (subHead != null ? subHead.hashCode() : 0);
        result = 31 * result + (description != null ? description.hashCode() : 0);
        result = 31 * result + Boolean.hashCode(isSalary);
        result = 31 * result + (bank != null ? bank.hashCode() : 0);
        return result;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Transaction t)) {
            return false;
        }
        return t.id == (id) &&
                t.date.equals(date) &&
                t.debit == debit &&
                t.credit == credit &&
                t.head.equals(head) &&
                t.subHead.equals(subHead) &&
                t.description.equals(description) &&
                t.isSalary == isSalary &&
                t.bank.equals(bank);
    }

    @Override
    public int compareTo(Transaction o) {
        return this.date.compareTo(o.date);
    }

    @Override
    public String toString() {
        return "Transaction {date=%s, debit=%d, credit=%d, head=%s, description=%s, isSalary=%b".formatted(getDateStr(), debit, credit, head, description, isSalary);
    }
}
