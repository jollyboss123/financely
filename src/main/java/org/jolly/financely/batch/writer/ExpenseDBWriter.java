package org.jolly.financely.batch.writer;

import org.jolly.financely.model.Expense;
import org.jolly.financely.repository.ExpenseRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.item.Chunk;
import org.springframework.batch.item.ItemWriter;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;

/**
 * @author jolly
 */
@Component
public class ExpenseDBWriter implements ItemWriter<Expense> {
    private static final Logger log = LoggerFactory.getLogger(ExpenseDBWriter.class);
    private final ExpenseRepository expenseRepository;

    public ExpenseDBWriter(ExpenseRepository expenseRepository) {
        this.expenseRepository = expenseRepository;
    }

    @Override
    public void write(@NonNull Chunk<? extends Expense> chunk) throws Exception {
        log.debug("writing expense data total: {} rows", chunk.getItems().size());
        expenseRepository.saveAll(chunk.getItems());
    }
}
