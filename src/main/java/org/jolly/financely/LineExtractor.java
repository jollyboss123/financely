package org.jolly.financely;

import java.util.List;

/**
 * @author jolly
 */
public interface LineExtractor {
    void dateRegex(String dateRegex);
    void linesToSkip(String[] linesToSkip);
    void startReadingText(String startReadingText);
    void endReadingText(String endReadingText);
    boolean extractLine(String pageContent, List<RawTransaction> items, String file);
}
