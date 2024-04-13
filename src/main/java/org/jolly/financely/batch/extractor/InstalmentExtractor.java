package org.jolly.financely.batch.extractor;

import jakarta.annotation.PostConstruct;
import org.jolly.financely.batch.extractor.DefaultFieldExtractor;
import org.springframework.stereotype.Component;

/**
 * @author jolly
 */
@Component
public class InstalmentExtractor extends DefaultFieldExtractor {
    @PostConstruct
    public void init() {
        super.setStringPatterns(new String[]{
                "\\d{2}\\/\\d{2}"
        });
    }
}
