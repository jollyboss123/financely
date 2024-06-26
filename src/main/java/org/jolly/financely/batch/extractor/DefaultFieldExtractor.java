package org.jolly.financely.batch.extractor;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * @author jolly
 */
public class DefaultFieldExtractor {
    private List<Pattern> patterns;
    private String[] stringPatterns;

    public DefaultFieldExtractor() {
        stringPatterns = new String[]{};
        patterns = new LinkedList<>();
    }

    public void setStringPatterns(String[] stringPatterns) {
        this.stringPatterns = stringPatterns;
        patterns = new LinkedList<>();
        for (String p : stringPatterns) {
            patterns.add(Pattern.compile(p));
        }
    }

    public String getField(String line) {
        Matcher m;
        for (Pattern p : patterns) {
            m = p.matcher(line);
            if (m.find()) {
                return m.group();
            }
        }
        return "";
    }
}
