package org.jolly.financely.exception;

/**
 * @author jolly
 */
public class PdfCloseException extends RuntimeException {
    public PdfCloseException(String message) {
        super(message);
    }
    public PdfCloseException(Throwable cause) {
        super(cause);
    }
}
