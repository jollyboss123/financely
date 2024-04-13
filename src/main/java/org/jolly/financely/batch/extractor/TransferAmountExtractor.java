package org.jolly.financely.batch.extractor;

import jakarta.annotation.PostConstruct;
import org.jolly.financely.batch.extractor.DefaultFieldExtractor;
import org.springframework.stereotype.Component;

/**
 * @author jolly
 */
@Component
public class TransferAmountExtractor extends DefaultFieldExtractor {
    @PostConstruct
    public void init() {
        super.setStringPatterns(new String[] {
                "(?<!\\d)\\d{1,3}(?:,\\d{3})+(?:\\.\\d{2,4})?",
                "\\d+\\.\\d+"
        });
    }
}
