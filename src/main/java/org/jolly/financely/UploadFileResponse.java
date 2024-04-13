package org.jolly.financely;

/**
 * @author jolly
 */
public class UploadFileResponse {
    private final String fileName;
    private final String fileType;
    private final long size;
    private final String message;
    private final String fileDownloadUri;
    private final Bank bank;

    UploadFileResponse(String fileName, String fileType, long size, String message, String fileDownloadUri, Bank bank) {
        this.fileName = fileName;
        this.fileType = fileType;
        this.size = size;
        this.message = message;
        this.fileDownloadUri = fileDownloadUri;
        this.bank = bank;
    }

    UploadFileResponse() {
        this(null, null, 0, null, null, null);
    }

    public static class Builder {
        private String fileName;
        private String fileType;
        private long size;
        private String message;
        private String fileDownloadUri;
        private Bank bank;

        public Builder fileName(String val) {
            fileName = val;
            return this;
        }

        public Builder fileType(String val) {
            fileType = val;
            return this;
        }

        public Builder size(long val) {
            size = val;
            return this;
        }

        public Builder message(String val) {
            message = val;
            return this;
        }

        public Builder fileDownloadUri(String val) {
            fileDownloadUri = val;
            return this;
        }

        public Builder uploadType(Bank val) {
            bank = val;
            return this;
        }

        public UploadFileResponse build() {
            return new UploadFileResponse(this);
        }
    }

    private UploadFileResponse(Builder builder) {
        this.fileName = builder.fileName;
        this.fileType = builder.fileType;
        this.size = builder.size;
        this.message = builder.message;
        this.fileDownloadUri = builder.fileDownloadUri;
        this.bank = builder.bank;
    }
}
