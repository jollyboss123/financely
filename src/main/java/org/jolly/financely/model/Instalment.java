package org.jolly.financely.model;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;

import java.util.Objects;

/**
 * @author jolly
 */
@Embeddable
public class Instalment {
    @Column(name = "instalment_number")
    private Integer number;
    @Column(name = "instalment_total")
    private Integer total;

    protected Instalment() {}

    protected Instalment(Integer number, Integer total) {
        this.number = number;
        this.total = total;
    }

    @Override
    public int hashCode() {
        int result = number.hashCode();
        result = 31 * result + (total != null ? total.hashCode() : 0);
        return result;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Instalment i)) {
            return false;
        }
        return Objects.equals(i.number, number) &&
                Objects.equals(i.total, total);
    }

    public static class Builder {
        private Integer number;
        private Integer total;

        public Builder number(Integer val) {
            number = val;
            return this;
        }

        public Builder total(Integer val) {
            total = val;
            return this;
        }

        public Instalment build() {
            return new Instalment(this);
        }
    }

    private Instalment(Builder builder) {
        this.number = builder.number;
        this.total = builder.total;
    }

    public Integer getNumber() {
        return number;
    }

    public void setNumber(Integer number) {
        this.number = number;
    }

    public Integer getTotal() {
        return total;
    }

    public void setTotal(Integer total) {
        this.total = total;
    }
}
