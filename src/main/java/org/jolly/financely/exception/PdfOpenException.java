package org.jolly.financely.exception;

/**
 * @author jolly
 */
public class PdfOpenException extends RuntimeException {
    public PdfOpenException(String message) {
        super(message);
    }
    public PdfOpenException(Throwable cause) {
        super(cause);
    }
}
