package org.jolly.financely;

import com.lowagie.text.pdf.PdfReader;
import com.lowagie.text.pdf.parser.PdfTextExtractor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.item.*;
import org.springframework.batch.item.file.ResourceAwareItemReaderItemStream;
import org.springframework.beans.factory.config.ConfigurableBeanFactory;
import org.springframework.context.annotation.Scope;
import org.springframework.core.io.Resource;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.LinkedList;
import java.util.List;

/**
 * @author jolly
 */
@Component
@Scope(value = ConfigurableBeanFactory.SCOPE_PROTOTYPE)
public class StatementPdfReader implements ResourceAwareItemReaderItemStream<RawTransaction> {
    private static final Logger log = LoggerFactory.getLogger(StatementPdfReader.class);
    private Resource resource;
    private String pdfPassword;
    private PdfReader pdfReader;
    private List<RawTransaction> items = new LinkedList<>();
    private int currentIndex = 0;
    private LineExtractor lineExtractor = new DefaultLineExtractor();

    public void setLineExtractor(LineExtractor lineExtractor) {
        this.lineExtractor = lineExtractor;
    }

    public void setPdfPassword(String password) {
        this.pdfPassword = password;
    }

    @Override
    public void setResource(@NonNull Resource resource) {
        this.resource = resource;
    }

    @Override
    public RawTransaction read() throws Exception, UnexpectedInputException, ParseException, NonTransientResourceException {
        if (currentIndex < items.size()) {
            return items.get(currentIndex++);
        }
        return null;
    }

    @Override
    public void open(ExecutionContext executionContext) throws ItemStreamException {
        log.info("started processing file: {}", resource);

        if (executionContext.containsKey("current.index")) {
            currentIndex = executionContext.getInt("current.index");
        } else {
            currentIndex = 0;
            try {
                readLines();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }

    @Override
    public void update(ExecutionContext executionContext) throws ItemStreamException {
        executionContext.putInt("current.index", currentIndex);
    }

    @Override
    public void close() throws ItemStreamException {
        log.debug("finished processing file: {}", resource);

        if (pdfReader != null && resource != null) {
            pdfReader.close();
        }
    }

    private void readLines() throws IOException {
        items = new LinkedList<>();
        if (pdfPassword != null) {
            pdfReader = new PdfReader(resource.getFile().getPath(), pdfPassword.getBytes());
        } else {
            pdfReader = new PdfReader(resource.getFile().getPath());
        }
        int pages = pdfReader.getNumberOfPages();

        PdfTextExtractor pdfTextExtractor = new PdfTextExtractor(pdfReader);

        for (int i = 2; i <= pages; i++) {
            String content = pdfTextExtractor.getTextFromPage(i, true);
            if (!lineExtractor.extractLine(content, items, resource.getFilename())) {
                break;
            }
        }
    }
}
