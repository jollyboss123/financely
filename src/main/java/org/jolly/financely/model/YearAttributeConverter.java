package org.jolly.financely.model;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

import java.time.Year;

/**
 * @author jolly
 */
@Converter(autoApply = true)
public class YearAttributeConverter implements AttributeConverter<Year, Short> {

    @Override
    public Short convertToDatabaseColumn(Year year) {
        if (year != null) {
            return (short) year.getValue();
        }
        return null;
    }

    @Override
    public Year convertToEntityAttribute(Short aShort) {
        if (aShort != null) {
            return Year.of(aShort);
        }
        return null;
    }
}
