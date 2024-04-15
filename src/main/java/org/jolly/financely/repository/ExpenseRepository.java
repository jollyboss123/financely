package org.jolly.financely.repository;

import org.jolly.financely.model.Expense;
import org.springframework.data.jpa.repository.JpaRepository;

/**
 * @author jolly
 */
public interface ExpenseRepository extends JpaRepository<Expense, Long> {
}
