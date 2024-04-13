package org.jolly.financely.model;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

import java.util.Currency;

/**
 * @author jolly
 */
@Converter(autoApply = true)
public class CurrencyConverter implements AttributeConverter<Currency, String> {
    @Override
    public String convertToDatabaseColumn(Currency currency) {
        if (currency != null) {
            return currency.getCurrencyCode();
        }
        return null;
    }

    @Override
    public Currency convertToEntityAttribute(String currencyCode) {
        if (currencyCode != null) {
            return Currency.getInstance(currencyCode);
        }
        return null;
    }
}
