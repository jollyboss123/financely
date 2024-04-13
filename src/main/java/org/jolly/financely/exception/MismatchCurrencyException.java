package org.jolly.financely.exception;

/**
 * Thrown when a set of <code>Money</code> objects do not have matching currencies.
 * <p>For example, adding together Euros and Dollars does not make any sense.
 *
 * @author jolly
 */
public class MismatchCurrencyException extends RuntimeException {
    public MismatchCurrencyException(String message) {
        super(message);
    }
}
