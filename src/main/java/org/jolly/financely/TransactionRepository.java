package org.jolly.financely;

import org.springframework.data.jpa.repository.JpaRepository;

/**
 * @author jolly
 */
public interface TransactionRepository extends JpaRepository<Transaction, Long> {
}
