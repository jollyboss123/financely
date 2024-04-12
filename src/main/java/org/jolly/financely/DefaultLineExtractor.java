package org.jolly.financely;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

/**
 * @author jolly
 */
public class DefaultLineExtractor implements LineExtractor {

    private static final Logger log = LoggerFactory.getLogger(DefaultLineExtractor.class);

    protected String startReadingText;
    protected String[] linesToSkip;
    protected String dateRegex = "^[0-9]{2}[a-zA-Z]{3}[0-9]{2}.*";
    protected String endReadingText;
    private RawTransaction rawTransaction;

    @Override
    public void dateRegex(String dateRegex) {
        this.dateRegex = dateRegex;
    }

    @Override
    public void linesToSkip(String[] linesToSkip) {
        this.linesToSkip = linesToSkip;
    }

    @Override
    public void startReadingText(String startReadingText) {
        this.startReadingText = startReadingText;
    }

    @Override
    public void endReadingText(String endReadingText) {
        this.endReadingText = endReadingText;
    }

    @Override
    public boolean extractLine(String pageContent, List<RawTransaction> items, String file) {
        boolean start = startReadingText == null;

        for (String line : pageContent.split("\\r?\\n")) {
            line = line.trim();
            log.debug("un-screened line: {}", line);
            if (shouldSkip(line)) {
                continue;
            }

            if (!start && line.matches(startReadingText)) {
                log.debug("starting line processing after this line: {}", startReadingText);
                start = true;
            }

            if (endReadingText != null && line.matches(endReadingText)) {
                log.debug("stopping line processing from this line: {}", endReadingText);
                return false;
            }

            if (start) {
                if (line.matches(dateRegex)) {
                    rawTransaction = new RawTransaction(file);
                    items.add(rawTransaction);
                }
                if (rawTransaction != null && !line.matches(startReadingText)) {
                    log.debug("read line: {}", line);
                    rawTransaction.getLines().add(line);
                }
            }
        }

        return true;
    }

    protected boolean shouldSkip(String line) {
        if (line.trim().isEmpty()) {
            return true;
        }

        if (linesToSkip != null) {
            for (String skipLine : linesToSkip) {
                if (line.matches(skipLine)) {
                    return true;
                }
            }
        }

        return false;
    }
}
