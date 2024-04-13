package org.jolly.financely.model;

import jakarta.persistence.*;
import org.jolly.financely.constant.Bank;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.Objects;

/**
 * @author jolly
 */
@Entity
@EntityListeners(AuditListener.class)
public class Transaction implements Comparable<Transaction>, Auditable {
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;
    @Embedded
    private Audit audit;
    @Temporal(TemporalType.DATE)
    private LocalDate date;
    private long debit;
    private long credit;
    private String head; // remove
    private String subHead; // remove
    private String description;
    private String file;
    private boolean isSalary;
    private boolean isInstalment;
    @Enumerated(EnumType.STRING)
    @Column(length = 8)
    private Bank bank;
    @Embedded
    private Instalment instalment;
    // from
    // to
    // currency
    private static final DateTimeFormatter DATE_FORMAT = DateTimeFormatter.ofPattern("yyyy-MM-dd");

    protected Transaction(long id, LocalDate date, long debit, long credit, String head, String subHead, String description, String file, boolean isSalary, boolean isInstalment, Bank bank, Instalment instalment) {
        this.id = id;
        this.date = date;
        this.debit = debit;
        this.credit = credit;
        this.head = head;
        this.subHead = subHead;
        this.description = description;
        this.file = file;
        this.isSalary = isSalary;
        this.bank = bank;
        this.isInstalment = isInstalment;
        this.instalment = instalment;
    }

    protected Transaction() {}

    public static class Builder {
        // required params
        private final long id;
        private final LocalDate date;
        private final Bank bank;
        private final String description;
        private final String file;

        // optional params
        private long debit = 0;
        private long credit = 0;
        private String head;
        private String subHead;
        private boolean isSalary = false;
        private boolean isInstalment = false;
        private Instalment instalment;

        public Builder(String file, long id, LocalDate date, Bank bank, String description) {
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

        public Builder isInstalment(boolean val) {
            isInstalment = val;
            return this;
        }

        public Builder instalment(Instalment val) {
            instalment = val;
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
        this.isInstalment = builder.isInstalment;
        this.instalment = builder.instalment;
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
        result = 31 * result + Boolean.hashCode(isInstalment);
        result = 31 * result + (bank != null ? bank.hashCode() : 0);
        result = 31 * result + (instalment != null ? instalment.hashCode() : 0);
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
        return Objects.equals(t.id, id) &&
                Objects.equals(t.date, date) &&
                t.debit == debit &&
                t.credit == credit &&
                Objects.equals(t.head, head) &&
                Objects.equals(t.subHead, subHead) &&
                Objects.equals(t.description, description) &&
                t.isSalary == isSalary &&
                Objects.equals(t.bank, bank);
    }

    @Override
    public int compareTo(Transaction o) {
        return this.date.compareTo(o.date);
    }

    @Override
    public String toString() {
        return "Transaction {date=%s, debit=%d, credit=%d, head=%s, description=%s, isSalary=%b".formatted(getDateStr(), debit, credit, head, description, isSalary);
    }

    public Long getId() {
        return id;
    }

    public Audit getAudit() {
        return audit;
    }

    public void setAudit(Audit audit) {
        this.audit = audit;
    }

    public String getFile() {
        return file;
    }

    public void setFile(String file) {
        this.file = file;
    }

    public LocalDate getDate() {
        return date;
    }

    public void setDate(LocalDate date) {
        this.date = date;
    }

    public long getCredit() {
        return credit;
    }

    public void setCredit(long credit) {
        this.credit = credit;
    }

    public long getDebit() {
        return debit;
    }

    public void setDebit(long debit) {
        this.debit = debit;
    }

    public Bank getBank() {
        return bank;
    }

    public void setBank(Bank bank) {
        this.bank = bank;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getHead() {
        return head;
    }

    public void setHead(String head) {
        this.head = head;
    }

    public String getSubHead() {
        return subHead;
    }

    public void setSubHead(String subHead) {
        this.subHead = subHead;
    }

    public boolean isSalary() {
        return isSalary;
    }

    public void setSalary(boolean salary) {
        isSalary = salary;
    }

    public Instalment getInstalment() {
        return instalment;
    }

    public void setInstalment(Instalment instalment) {
        this.instalment = instalment;
    }
}
