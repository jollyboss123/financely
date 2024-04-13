package org.jolly.financely.money;

import jakarta.persistence.Column;
import jakarta.persistence.Embeddable;
import jakarta.persistence.Transient;
import jakarta.validation.constraints.NotNull;
import org.springframework.lang.NonNull;

import java.io.*;
import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.Collection;
import java.util.Currency;
import java.util.Objects;

/**
 * @author jolly
 */
@Embeddable
public final class Money implements Comparable<Money>, Serializable {
    @Serial
    private static final long serialVersionUID = -3481031843124881066L;

    @Column(name = "amount")
    private BigDecimal amount;
    @Column(name = "currency", length = 10)
    private final Currency currency;
    @Transient
    private final RoundingMode rounding;
    private static Currency DEFAULT_CURRENCY;
    private static RoundingMode DEFAULT_ROUNDING;

    protected Money(BigDecimal amount, Currency currency, RoundingMode rounding, boolean autoRound) {
        this.amount = amount;
        this.currency = currency;
        this.rounding = rounding;
        validateState(autoRound);
    }

    protected Money(BigDecimal amount, Currency currency, RoundingMode rounding) {
        this(amount, currency, rounding, false);
    }

    private Money(BigDecimal amount, Currency currency, boolean autoRound) {
        this(amount, currency, DEFAULT_ROUNDING, autoRound);
    }

    private Money(BigDecimal amount, Currency currency) {
        this(amount, currency, DEFAULT_ROUNDING);
    }

    private Money(BigDecimal amount, boolean autoRound) {
        this(amount, DEFAULT_CURRENCY, DEFAULT_ROUNDING, autoRound);
    }

    private Money(BigDecimal amount) {
        this(amount, DEFAULT_CURRENCY, DEFAULT_ROUNDING);
    }

    protected Money() {
        this(BigDecimal.ZERO, DEFAULT_CURRENCY, DEFAULT_ROUNDING);
    }

    public static void init(@NonNull Currency defaultCurrency, @NonNull RoundingMode defaultRounding) {
        DEFAULT_CURRENCY = Objects.requireNonNull(defaultCurrency);
        DEFAULT_ROUNDING = Objects.requireNonNull(defaultRounding);
    }

    public static Money of(@NonNull BigDecimal amount, @NonNull Currency currency, @NonNull RoundingMode rounding, boolean autoRound) {
        return new Money(amount, currency, rounding, autoRound);
    }

    public static Money of(@NonNull BigDecimal amount, @NonNull Currency currency, boolean autoRound) {
        return new Money(amount, currency, autoRound);
    }

    public static Money of(@NonNull BigDecimal amount, boolean autoRound) {
        return new Money(amount, autoRound);
    }

    public static Money of(@NonNull BigDecimal amount, @NonNull Currency currency, @NonNull RoundingMode rounding) {
        return new Money(amount, currency, rounding);
    }

    public static Money of(@NonNull BigDecimal amount, @NonNull Currency currency) {
        return new Money(amount, currency);
    }

    public static Money of(@NonNull BigDecimal amount) {
        return new Money(amount);
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public Currency getCurrency() {
        return currency;
    }

    public RoundingMode getRoundingStyle() {
        return rounding;
    }

    public boolean isSameCurrency(Money that) {
        boolean result = false;
        if (that != null) {
            result = this.currency.equals(that.currency);
        }
        return result;
    }

    public Money plus(Money that) {
        checkCurrenciesMatch(that);
        return new Money(amount.add(that.amount), currency, rounding);
    }

    public Money minus(Money that) {
        checkCurrenciesMatch(that);
        return new Money(amount.subtract(that.amount), currency, rounding);
    }

    public static Money sum(Collection<Money> monies) {
        Money sum = new Money(BigDecimal.ZERO);
        for (Money m : monies) {
            sum = sum.plus(m);
        }
        return sum;
    }

    public Money times(int factor) {
        BigDecimal aFactor = new BigDecimal(factor);
        BigDecimal newAmount = amount.multiply(aFactor);
        return new Money(newAmount, currency, rounding);
    }

    public Money times(double factor) {
        BigDecimal aFactor = asBigDecimal(factor);
        BigDecimal newAmount = amount.multiply(aFactor)
                .setScale(getNumDecimalsForCurrency(), rounding);
        return new Money(newAmount, currency, rounding);
    }

    public Money div(int divisor) {
        BigDecimal aDivisor = new BigDecimal(divisor);
        BigDecimal newAmount = amount.divide(aDivisor, rounding);
        return new Money(newAmount, currency, rounding);
    }

    public Money div(double divisor) {
        BigDecimal aDivisor = asBigDecimal(divisor);
        BigDecimal newAmount = amount.divide(aDivisor, rounding);
        return new Money(newAmount, currency, rounding);
    }

    public boolean eq(Money that) {
        checkCurrenciesMatch(that);
        return compareAmount(that) == 0;
    }

    public boolean gt(Money that) {
        checkCurrenciesMatch(that);
        return compareAmount(that) > 0;
    }

    public boolean gteq(Money that) {
        checkCurrenciesMatch(that);
        return compareAmount(that) >= 0;
    }

    public boolean lt(Money that) {
        checkCurrenciesMatch(that);
        return compareAmount(that) < 0;
    }

    public boolean lteq(Money that) {
        checkCurrenciesMatch(that);
        return compareAmount(that) <= 0;
    }

    private Object[] getSigFields() {
        return new Object[]{amount, currency, rounding};
    }

    @Override
    public String toString() {
        return amount.toPlainString() + " " + currency.getSymbol();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Money m)) {
            return false;
        }
        for (int i = 0; i < this.getSigFields().length; i++) {
            if (!Objects.equals(this.getSigFields()[i], m.getSigFields()[i])) {
                return false;
            }
        }
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.getSigFields());
    }

    @Override
    public int compareTo(@NotNull Money o) {
        final int equal = 0;
        if (this == o) {
            return equal;
        }
        int comparison = this.amount.compareTo(o.amount);
        if (comparison != equal) {
            return comparison;
        }
        comparison = this.currency.getCurrencyCode().compareTo(o.currency.getCurrencyCode());
        if (comparison != equal) {
            return comparison;
        }
        comparison = this.rounding.compareTo(o.rounding);
        if (comparison != equal) {
            return comparison;
        }
        return equal;
    }

    @Serial
    private void readObject(ObjectInputStream in) throws IOException, ClassNotFoundException {
        in.defaultReadObject();
        // defensive copy for mutable field since BigDecimal is non-final
        amount = new BigDecimal(amount.toPlainString());
        validateState(false);
    }

    @Serial
    private void writeObject(ObjectOutputStream out) throws IOException {
        out.defaultWriteObject();
    }

    private void validateState(boolean autoRound) {
        if (amount == null) {
            throw new IllegalArgumentException("amount cannot be null");
        }

        if (currency == null) {
            throw new IllegalArgumentException("currency cannot be null");
        }

        if (amount.scale() > getNumDecimalsForCurrency()) {
            if (!autoRound) {
                throw new IllegalArgumentException("number of decimals is " + amount.scale() + ", but currency only takes " +
                        getNumDecimalsForCurrency() + " decimals.");
            }
            amount = amount.setScale(getNumDecimalsForCurrency(), rounding);
        }
    }

    private int getNumDecimalsForCurrency() {
        return currency.getDefaultFractionDigits();
    }

    private void checkCurrenciesMatch(Money that) {
        if (!this.currency.equals(that.getCurrency())) {
            throw new MismatchCurrencyException(that.getCurrency() + "doesn't match the expected currency: " + this.currency);
        }
    }

    /**
     * Ignores scale: 0 same as 0.00.
     */
    private int compareAmount(Money that) {
        return this.amount.compareTo(that.amount);
    }

    private BigDecimal asBigDecimal(double aDouble) {
        String s = Double.toString(aDouble);
        return new BigDecimal(s);
    }
}
