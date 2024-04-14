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
public class MaybankStatementBatchJob {
    private static final Logger log = LoggerFactory.getLogger(MaybankStatementBatchJob.class);
    @Value("file:${file.path.mbb}")
    private Resource[] resources;
    private static final String JOB_NAME = "MBBAccount.ETL.Job";
    private static final String PROCESSOR_TASK_NAME = "MBBAccount.ETL.Job.file.load";

    @Bean
    public Job mbbBankJob(JobRepository jobRepository,
                          PlatformTransactionManager transactionManager,
                          ItemReader<RawTransaction> mbbItemsReader,
                          ItemProcessor<RawTransaction, Transaction> mbbItemProcessor,
                          ItemWriter<Transaction> bankAccountDBWriter) {
        Step step = new StepBuilder(PROCESSOR_TASK_NAME, jobRepository)
                .<RawTransaction, Transaction>chunk(100, transactionManager)
                .reader(mbbItemsReader)
                .processor(mbbItemProcessor)
                .writer(bankAccountDBWriter)
                .build();

        return new JobBuilder(JOB_NAME, jobRepository)
                .incrementer(new RunIdIncrementer())
                .start(step)
                .build();
    }

    @Bean
    public MultiResourceItemReader<RawTransaction> mbbItemsReader(StatementPdfReader mbbItemReader) {
        MultiResourceItemReader<RawTransaction> reader = new MultiResourceItemReader<>();
        reader.setResources(resources);
        reader.setStrict(false);
        reader.setDelegate(mbbItemReader);
        return reader;
    }

    @Bean
    public StatementPdfReader mbbItemReader(@Qualifier("StatementPdfReader") StatementPdfReader flatFileItemReader) {
        LineExtractor defaultLineExtractor = new DefaultLineExtractor();
        defaultLineExtractor.dateRegex("^[0-9]{2}\\/[0-9]{2}\\/[0-9]{4}.*");
        defaultLineExtractor.startReadingText(".*ENTRY DATE.*");
        defaultLineExtractor.endReadingText(".*ENDING BALANCE.*");
        defaultLineExtractor.linesToSkip(new String[]{
                "^Perhation / Note.*"
        });
        flatFileItemReader.setLineExtractor(defaultLineExtractor);
        flatFileItemReader.setPdfPassword(""); //TODO: use bouncy castle
        return flatFileItemReader;
    }

    @Bean
    public BankAccountProcessor mbbItemProcessor(@Qualifier("BankAccountProcessor") BankAccountProcessor itemProcessor) {
        itemProcessor.setDateTimeFormatter(DateTimeFormatter.ofPattern("dd/MM/yyyy"));
        itemProcessor.setDateLengths(new BankAccountProcessor.DateLength(10, null));
        itemProcessor.setCreditTransfer(new String[]{
                "(?<!\\d)\\d{1,3}(?:,\\d{3})+(?:\\.\\d{2})?\\+",
                "\\d+\\.\\d+\\+"
        });
        return itemProcessor;
    }
}
