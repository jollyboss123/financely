package org.jolly.financely.batch.writer;

import org.jolly.financely.model.Transaction;
import org.jolly.financely.repository.TransactionRepository;
import org.springframework.batch.item.Chunk;
import org.springframework.batch.item.ItemWriter;
import org.springframework.lang.NonNull;
import org.springframework.stereotype.Component;

/**
 * @author jolly
 */
@Component
public class BankAccountDBWriter implements ItemWriter<Transaction> {
    private final TransactionRepository transactionRepository;

    public BankAccountDBWriter(TransactionRepository transactionRepository) {
        this.transactionRepository = transactionRepository;
    }

    @Override
    public void write(@NonNull Chunk<? extends Transaction> chunk) throws Exception {
        transactionRepository.saveAll(chunk.getItems());
    }
}
