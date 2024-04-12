package org.jolly.financely;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.job.builder.JobBuilder;
import org.springframework.batch.core.launch.support.RunIdIncrementer;
import org.springframework.batch.core.repository.JobRepository;
import org.springframework.batch.core.step.builder.StepBuilder;
import org.springframework.batch.item.Chunk;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.batch.item.ItemReader;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.file.MultiResourceItemReader;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.Resource;
import org.springframework.lang.NonNull;
import org.springframework.transaction.PlatformTransactionManager;

import java.time.format.DateTimeFormatterBuilder;
import java.time.temporal.ChronoField;
import java.util.Locale;

/**
 * @author jolly
 */
@Configuration
public class UOBStatementBatchJob {
    private static final Logger log = LoggerFactory.getLogger(UOBStatementBatchJob.class);
    @Value("file:${file.path.uob}")
    private Resource[] resources;
    private static final String JOB_NAME = "UOBAccount.ETL.Job";
    private static final String PROCESSOR_TASK_NAME = "UOBAccount.ETL.Job.file.load";
    private static final String ARCHIVE_TASK_NAME = "UOBAccount.ETL.Job.file.archive";

    @Bean
    public Job uobBankJob(JobRepository jobRepository,
                          PlatformTransactionManager transactionManager,
                          ItemReader<RawTransaction> uobItemsReader,
                          ItemProcessor<RawTransaction, Transaction> uobItemProcessor,
                          ItemWriter<Transaction> uobItemWriter
                          ) {
        Step step = new StepBuilder(PROCESSOR_TASK_NAME, jobRepository)
                .<RawTransaction, Transaction>chunk(100, transactionManager)
                .reader(uobItemsReader)
                .processor(uobItemProcessor)
                .writer(uobItemWriter)
                .build();

        return new JobBuilder(JOB_NAME, jobRepository)
                .incrementer(new RunIdIncrementer())
                .start(step)
                .build();
    }

    @Bean
    public BankAccountProcessor uobItemProcessor(BankAccountProcessor itemProcessor) {
        itemProcessor.setDateTimeFormatter(new DateTimeFormatterBuilder()
                .parseCaseInsensitive()
                .appendPattern("dd MMM")
                .parseDefaulting(ChronoField.YEAR, 2024) //TODO: parse year from statement automatically
                .toFormatter(Locale.getDefault()));
        itemProcessor.setDateLen(6);
        return itemProcessor;
    }

    @Bean
    public ItemWriter<Transaction> uobItemWriter() {
        return new ItemWriter<Transaction>() {
            @Override
            public void write(@NonNull Chunk<? extends Transaction> chunk) throws Exception {
                for (Transaction c : chunk.getItems()) {
                    if (log.isInfoEnabled()) {
                        log.info(c.toString());
                    }
                }
            }
        };
    }

    @Bean
    public MultiResourceItemReader<RawTransaction> uobItemsReader(StatementPdfReader uobItemReader) {
        MultiResourceItemReader<RawTransaction> reader = new MultiResourceItemReader<>();
        reader.setResources(resources);
        reader.setStrict(false);
        reader.setDelegate(uobItemReader);
        return reader;
    }

    @Bean
    public StatementPdfReader uobItemReader(StatementPdfReader flatFileItemReader) {
        LineExtractor defaultLineExtractor = new DefaultLineExtractor();
        defaultLineExtractor.dateRegex("^[0-9]{2} [a-zA-Z]{3}.*");
        defaultLineExtractor.startReadingText("Transaction Date");
        defaultLineExtractor.endReadingText("END OF STATEMENT*");
        defaultLineExtractor.linesToSkip(
                new String[]{

                }
        );
        flatFileItemReader.setLineExtractor(defaultLineExtractor);
        return flatFileItemReader;
    }
}
