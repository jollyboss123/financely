package org.jolly.financely.model;

import java.util.LinkedList;
import java.util.List;

/**
 * @author jolly
 */
public class RawTransaction {
    private final List<String> lines;
    private final String file;

    public RawTransaction(String file) {
        this(new LinkedList<>(), file);
    }

    public RawTransaction(List<String> lines, String file) {
        this.lines = lines;
        this.file = file;
    }

    public String getMergedLines(int from) {
        final String str = String.join(", ", lines);
        final int length = Math.min(str.length(), 252);

        return str.substring(from, length);
    }

    public List<String> getLines() {
        return lines;
    }

    public String getFile() {
        return file;
    }

    @Override
    public String toString() {
        return "RawTransaction {file=%s, content=%s}".formatted(file, lines.toString());
    }
}
