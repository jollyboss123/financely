package org.jolly.financely.repository;

import org.jolly.financely.model.Transaction;
import org.springframework.data.jpa.repository.JpaRepository;

/**
 * @author jolly
 */
public interface TransactionRepository extends JpaRepository<Transaction, Long> {
}
