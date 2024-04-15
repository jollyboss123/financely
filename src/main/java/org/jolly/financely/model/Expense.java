package org.jolly.financely.model;

import jakarta.persistence.*;
import org.springframework.lang.NonNull;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.Month;
import java.time.Year;
import java.util.Objects;

/**
 * @author jolly
 */
@Entity
@EntityListeners(AuditListener.class)
public class Expense implements Comparable<Expense>, Auditable {
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;
    @Embedded
    private Audit audit;
    @Temporal(TemporalType.TIMESTAMP)
    private LocalDateTime date;
    @Embedded
    private Money price;
    private String description;
    private String category;
    @Column(columnDefinition = "smallint")
    @Convert(converter = YearAttributeConverter.class)
    private Year year;
    @Column(columnDefinition = "smallint")
    @Enumerated
    private Month month;

    protected Expense(Long id, LocalDateTime date, Money price, String description, String category, Year year, Month month) {
        this.id = id;
        this.date = date;
        this.price = price;
        this.description = description;
        this.category = category;
        this.year = year;
        this.month = month;
    }

    protected Expense() {}

    public static class Builder {
        // required
        private final LocalDateTime date;
        private final Money price;
        private final String category;

        // optional
        private String description;
        private Year year;
        private Month month;

        public Builder(@NonNull LocalDateTime date, @NonNull Money price, @NonNull String category) {
            this.date = date;
            this.price = price;
            this.category = category;
        }

        public Builder description(String val) {
            this.description = val;
            return this;
        }

        public Builder year(Year val) {
            this.year = val;
            return this;
        }

        public Builder month(Month val) {
            this.month = val;
            return this;
        }

        public Expense build() {
            if (year == null) {
                year = Year.of(date.getYear());
            }
            if (month == null) {
                month = Month.of(date.getMonthValue());
            }
            return new Expense(this);
        }
    }

    private Expense(Builder builder) {
        this.date = builder.date;
        this.price = builder.price;
        this.category = builder.category;
        this.description = builder.description;
        this.year = builder.year;
        this.month = builder.month;
    }

    @Override
    public int hashCode() {
        int result = Long.hashCode(id);
        result = 31 * result + date.hashCode();
        result = 31 * result + price.hashCode();
        result = 31 * result + category.hashCode();
        result = 31 * result + (description != null ? description.hashCode() : 0);
        result = 31 * result + (year != null ? year.hashCode() : 0);
        result = 31 * result + (month != null ? month.hashCode() : 0);
        return result;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Expense e)) {
            return false;
        }
        return Objects.equals(e.id, id) &&
                Objects.equals(e.date, date) &&
                Objects.equals(e.price, price) &&
                Objects.equals(e.category, category) &&
                Objects.equals(e.description, description) &&
                Objects.equals(e.month, month) &&
                Objects.equals(e.year, year);
    }

    @Override
    public int compareTo(@NonNull Expense o) {
        return this.date.compareTo(o.date);
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

    public LocalDateTime getDate() {
        return date;
    }

    public void setDate(LocalDateTime date) {
        this.date = date;
    }

    public Money getPrice() {
        return price;
    }

    public void setPrice(Money price) {
        this.price = price;
    }

    public String getCategory() {
        return category;
    }

    public void setCategory(String category) {
        this.category = category;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Year getYear() {
        return year;
    }

    public void setYear(Year year) {
        this.year = year;
    }

    public Month getMonth() {
        return month;
    }

    public void setMonth(Month month) {
        this.month = month;
    }
}
