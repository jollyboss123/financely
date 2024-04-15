package org.jolly.financely.batch.processor;

import org.jolly.financely.model.Expense;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.beans.factory.config.ConfigurableBeanFactory;
import org.springframework.context.annotation.Scope;
import org.springframework.stereotype.Component;

/**
 * @author jolly
 */
@Component
@Scope(value = ConfigurableBeanFactory.SCOPE_PROTOTYPE)
public class DefaultExpenseProcessor implements ItemProcessor<Expense, Expense> {
    private static final Logger log = LoggerFactory.getLogger(DefaultExpenseProcessor.class);
    private String[] creditTransfer;
    private String[] itemsToSkip;

    public void setCreditTransfer(String[] creditTransfer) {
        this.creditTransfer = creditTransfer;
    }

    public void setItemsToSkip(String[] itemsToSkip) {
        this.itemsToSkip = itemsToSkip;
    }

    @Override
    public Expense process(Expense item) throws Exception {
        if (shouldSkip(item.getDescription()) || shouldSkip(item.getCategory())) {
            return null;
        }

        return item;
    }

    private boolean shouldSkip(String s) {
        if (s.trim().isEmpty()) {
            return true;
        }

        if (itemsToSkip != null) {
            for (String item : itemsToSkip) {
                if (s.matches(item)) {
                    return true;
                }
            }
        }

        return false;
    }
}
