package org.jolly.financely;

import org.jolly.financely.annotation.LogExecutionTime;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.core.*;
import org.springframework.batch.core.launch.JobLauncher;
import org.springframework.batch.core.repository.JobExecutionAlreadyRunningException;
import org.springframework.batch.core.repository.JobInstanceAlreadyCompleteException;
import org.springframework.batch.core.repository.JobRestartException;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

/**
 * @author jolly
 */
@RestController
@RequestMapping("/expense")
public class ExpenseController {
    private static final Logger log = LoggerFactory.getLogger(ExpenseController.class);
    private final JobLauncher jobLauncher;
    private final Job job;

    public ExpenseController(JobLauncher jobLauncher, @Qualifier("expenseJob") Job job) {
        this.jobLauncher = jobLauncher;
        this.job = job;
    }

    @LogExecutionTime
    @GetMapping("/load")
    public BatchStatus load() throws JobInstanceAlreadyCompleteException, JobExecutionAlreadyRunningException, JobParametersInvalidException, JobRestartException {
        Map<String, JobParameter<?>> parameters = new HashMap<>();
        parameters.put("time", new JobParameter<>(System.currentTimeMillis(), Long.class));
        JobParameters jobParameters = new JobParameters(parameters);

        JobExecution jobExecution = jobLauncher.run(job, jobParameters);
        log.debug("batch status: {}", jobExecution.getStatus());

        while (jobExecution.isRunning()) {
            log.debug("batch is running");
            log.debug("...");
        }

        return jobExecution.getStatus();
    }
}
