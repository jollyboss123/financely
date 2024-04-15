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
    @Embedded
    @AttributeOverride(name = "amount", column = @Column(name = "credit_amount"))
    @AttributeOverride(name = "currency", column = @Column(name = "credit_currency", length = 10))
    private Money credit;
    @Embedded
    @AttributeOverride(name = "amount", column = @Column(name = "debit_amount"))
    @AttributeOverride(name = "currency", column = @Column(name = "debit_currency", length = 10))
    private Money debit;
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
    private static final DateTimeFormatter DATE_FORMAT = DateTimeFormatter.ofPattern("yyyy-MM-dd");

    protected Transaction(Long id, LocalDate date, Money debit, Money credit, String description, String file, boolean isSalary, boolean isInstalment, Bank bank, Instalment instalment) {
        this.id = id;
        this.date = date;
        this.debit = debit;
        this.credit = credit;
        this.description = description;
        this.file = file;
        this.isSalary = isSalary;
        this.bank = bank;
        this.isInstalment = isInstalment;
        this.instalment = instalment;
    }

    protected Transaction() {}

    public static class Builder {
        // required
        private final LocalDate date;
        private final Bank bank;
        private final String description;
        private final String file;

        // optional
        private Money debit;
        private Money credit;
        private boolean isSalary = false;
        private boolean isInstalment = false;
        private Instalment instalment;

        public Builder(String file, LocalDate date, Bank bank, String description) {
            this.file = file;
            this.date = date;
            this.bank = bank;
            this.description = description;
        }

        public Builder debit(Money val) {
            debit = val;
            return this;
        }

        public Builder credit(Money val) {
            credit = val;
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
        this.bank = builder.bank;
        this.description = builder.description;
        this.date = builder.date;
        this.debit = builder.debit;
        this.credit = builder.credit;
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
        result = 31 * result + (debit != null ? debit.hashCode() : 0);
        result = 31 * result + (credit != null ? credit.hashCode() : 0);
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
                Objects.equals(t.debit, debit) &&
                Objects.equals(t.credit, credit) &&
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
        return "Transaction {date=%s, debit=%s, credit=%s, description=%s, isSalary=%b".formatted(getDateStr(), debit.toString(), credit.toString(), description, isSalary);
    }

    public Long getId() {
        return id;
    }

    @Override
    public Audit getAudit() {
        return audit;
    }

    @Override
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

    public Money getCredit() {
        return credit;
    }

    public void setCredit(Money credit) {
        this.credit = credit;
    }

    public Money getDebit() {
        return debit;
    }

    public void setDebit(Money debit) {
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
