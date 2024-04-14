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

import java.time.format.DateTimeFormatterBuilder;
import java.time.temporal.ChronoField;
import java.util.Locale;

/**
 * @author jolly
 */
@Configuration
public class GXStatementBatchJob {
    private static final Logger log = LoggerFactory.getLogger(GXStatementBatchJob.class);
    @Value("file:${file.path.gx}")
    private Resource[] resources;
    private static final String JOB_NAME = "GXAccount.ETL.Job";
    private static final String PROCESSOR_TASK_NAME = "GXAccount.ETL.Job.file.load";

    @Bean
    public Job gxBankJob(JobRepository jobRepository,
                         PlatformTransactionManager transactionManager,
                         ItemReader<RawTransaction> gxItemsReader,
                         ItemProcessor<RawTransaction, Transaction> gxItemProcessor,
                         ItemWriter<Transaction> bankAccountDBWriter) {
        Step step = new StepBuilder(PROCESSOR_TASK_NAME, jobRepository)
                .<RawTransaction, Transaction>chunk(100, transactionManager)
                .reader(gxItemsReader)
                .processor(gxItemProcessor)
                .writer(bankAccountDBWriter)
                .build();

        return new JobBuilder(JOB_NAME, jobRepository)
                .incrementer(new RunIdIncrementer())
                .start(step)
                .build();
    }

    @Bean
    public MultiResourceItemReader<RawTransaction> gxItemsReader(StatementPdfReader gxItemReader) {
        MultiResourceItemReader<RawTransaction> reader = new MultiResourceItemReader<>();
        reader.setResources(resources);
        reader.setStrict(false);
        reader.setDelegate(gxItemReader);
        return reader;
    }

    @Bean
    public StatementPdfReader gxItemReader(@Qualifier("StatementPdfReader") StatementPdfReader flatFileItemReader) {
        LineExtractor defaultLineExtractor = new DefaultLineExtractor();
        defaultLineExtractor.dateRegex("^[0-9]{1,2} [a-zA-Z]{3}.*");
        defaultLineExtractor.startReadingText(".*Transaction description.*");
        defaultLineExtractor.linesToSkip(new String[]{
                ".*GX Bank Berhad formerly known as.*",
                "^Page .*"
        });
        flatFileItemReader.setLineExtractor(defaultLineExtractor);
        return flatFileItemReader;
    }

    @Bean
    public BankAccountProcessor gxItemProcessor(@Qualifier("BankAccountProcessor") BankAccountProcessor itemProcessor) {
        itemProcessor.setDateTimeFormatter(new DateTimeFormatterBuilder()
                .parseCaseInsensitive()
                .appendPattern("d MMM")
                .parseDefaulting(ChronoField.YEAR, 2024) //TODO: parse year from statement automatically
                .toFormatter(Locale.getDefault()));
        itemProcessor.setDateLengths(new BankAccountProcessor.DateLength(5, 6));
        itemProcessor.setCreditTransfer(new String[]{
                ".*Interest earned.*",
                "\\+(?<!\\d)\\d{1,3}(?:,\\d{3})+(?:\\.\\d{2})?",
                "\\+\\d+\\.\\d+"
        });
        return itemProcessor;
    }
}
