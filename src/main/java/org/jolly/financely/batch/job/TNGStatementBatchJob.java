package org.jolly.financely.batch.job;

import org.jolly.financely.batch.extractor.DefaultLineExtractor;
import org.jolly.financely.batch.extractor.LineExtractor;
import org.jolly.financely.batch.processor.BankAccountProcessor;
import org.jolly.financely.batch.reader.StatementPdfReader;
import org.jolly.financely.model.RawTransaction;
import org.jolly.financely.model.Transaction;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.job.builder.JobBuilder;
import org.springframework.batch.core.launch.support.RunIdIncrementer;
import org.springframework.batch.core.repository.JobRepository;
import org.springframework.batch.core.step.builder.StepBuilder;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.batch.item.ItemReader;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.file.MultiResourceItemReader;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.Resource;
import org.springframework.transaction.PlatformTransactionManager;

import java.time.format.DateTimeFormatter;

/**
 * @author jolly
 */
@Configuration
public class TNGStatementBatchJob {
    private static final Logger log = LoggerFactory.getLogger(TNGStatementBatchJob.class);
    @Value("file:${file.path.tng}")
    private Resource[] resources;
    private static final String JOB_NAME = "TNGAccount.ETL.Job";
    private static final String PROCESSOR_TASK_NAME = "TNGAccount.ETL.Job.file.load";

    @Bean
    public Job tngBankJob(JobRepository jobRepository,
                           PlatformTransactionManager transactionManager,
                           ItemReader<RawTransaction> tngItemsReader,
                           ItemProcessor<RawTransaction, Transaction> tngItemProcessor,
                           ItemWriter<Transaction> bankAccountDBWriter) {
        Step step = new StepBuilder(PROCESSOR_TASK_NAME, jobRepository)
                .<RawTransaction, Transaction>chunk(100, transactionManager)
                .reader(tngItemsReader)
                .processor(tngItemProcessor)
                .writer(bankAccountDBWriter)
                .build();

        return new JobBuilder(JOB_NAME, jobRepository)
                .incrementer(new RunIdIncrementer())
                .start(step)
                .build();
    }

    @Bean
    public MultiResourceItemReader<RawTransaction> tngItemsReader(StatementPdfReader tngItemReader) {
        MultiResourceItemReader<RawTransaction> reader = new MultiResourceItemReader<>();
        reader.setResources(resources);
        reader.setStrict(false);
        reader.setDelegate(tngItemReader);
        return reader;
    }

    @Bean
    public StatementPdfReader tngItemReader(@Qualifier("StatementPdfReader") StatementPdfReader flatFileItemReader) {
        LineExtractor defaultLineExtractor = new DefaultLineExtractor();
        defaultLineExtractor.dateRegex("^[0-9]{1,2}\\/[0-9]{1,2}\\/[0-9]{4}.*");
        defaultLineExtractor.linesToSkip(new String[]{
                "^\\*This is a system generated email\\..*"
        });
        flatFileItemReader.setLineExtractor(defaultLineExtractor);
        return flatFileItemReader;
    }

    @Bean
    public BankAccountProcessor tngItemProcessor(@Qualifier("BankAccountProcessor") BankAccountProcessor itemProcessor) {
        itemProcessor.setDateTimeFormatter(DateTimeFormatter.ofPattern("d/M/yyyy"));
        itemProcessor.setDateLengths(new BankAccountProcessor.DateLength(8, 10));
        itemProcessor.setCreditTransfer(new String[]{
                ".*DUITNOW_RECEIVEFROM.*",
                ".*Receive from Wallet.*",
                ".*Daily Earnings.*"
        });
        return itemProcessor;
    }
}
