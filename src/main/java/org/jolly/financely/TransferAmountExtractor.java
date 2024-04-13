package org.jolly.financely;

import jakarta.annotation.PostConstruct;
import org.springframework.stereotype.Component;

/**
 * @author jolly
 */
@Component
public class TransferAmountExtractor extends DefaultFieldExtractor {
    @PostConstruct
    public void init() {
        super.setStringPatterns(new String[] {
                "(?<!\\d)\\d{1,3}(?:,\\d{3})+(?:\\.\\d{2})?",
                "\\d+\\.\\d+"
        });
    }
}
