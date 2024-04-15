package org.jolly.financely.batch.job;

import org.jolly.financely.batch.extractor.DefaultLineExtractor;
import org.jolly.financely.batch.extractor.LineExtractor;
import org.jolly.financely.batch.processor.BankAccountProcessor;
import org.jolly.financely.batch.reader.PdfReader;
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
public class CIMBStatementBatchJob {
    private static final Logger log = LoggerFactory.getLogger(CIMBStatementBatchJob.class);
    @Value("file:${file.path.cimb}")
    private Resource[] resources;
    private static final String JOB_NAME = "CIMBAccount.ETL.Job";
    private static final String PROCESSOR_TASK_NAME = "CIMBAccount.ETL.Job.file.load";

    @Bean
    public Job cimbBankJob(JobRepository jobRepository,
                           PlatformTransactionManager transactionManager,
                           ItemReader<RawTransaction> cimbItemsReader,
                           ItemProcessor<RawTransaction, Transaction> cimbItemProcessor,
                           ItemWriter<Transaction> bankAccountDBWriter) {
        Step step = new StepBuilder(PROCESSOR_TASK_NAME, jobRepository)
                .<RawTransaction, Transaction>chunk(100, transactionManager)
                .reader(cimbItemsReader)
                .processor(cimbItemProcessor)
                .writer(bankAccountDBWriter)
                .build();

        return new JobBuilder(JOB_NAME, jobRepository)
                .incrementer(new RunIdIncrementer())
                .start(step)
                .build();
    }

    @Bean
    public MultiResourceItemReader<RawTransaction> cimbItemsReader(PdfReader cimbItemReader) {
        MultiResourceItemReader<RawTransaction> reader = new MultiResourceItemReader<>();
        reader.setResources(resources);
        reader.setStrict(false);
        reader.setDelegate(cimbItemReader);
        return reader;
    }

    @Bean
    public PdfReader cimbItemReader(@Qualifier("pdfReader") PdfReader flatFileItemReader) {
        LineExtractor defaultLineExtractor = new DefaultLineExtractor();
        defaultLineExtractor.dateRegex("^[0-9]{2}\\/[0-9]{2}\\/[0-9]{4}.*");
        defaultLineExtractor.startReadingText(".*Ref No.*");
        defaultLineExtractor.endReadingText(".*End of Statement.*");
        defaultLineExtractor.linesToSkip(new String[]{
                "^Important Notice.*",
                "^Effective 8 November 2021.*",
                "^The Bank must be informed of any error.*",
                "^You can transfer funds, enquire balances.*"
        });
        flatFileItemReader.setLineExtractor(defaultLineExtractor);
        return flatFileItemReader;
    }

    @Bean
    public BankAccountProcessor cimbItemProcessor(@Qualifier("bankAccountProcessor") BankAccountProcessor itemProcessor) {
        itemProcessor.setDateTimeFormatter(DateTimeFormatter.ofPattern("dd/MM/yyyy"));
        itemProcessor.setDateLengths(new BankAccountProcessor.DateLength(10, null));
        itemProcessor.setCreditTransfer(new String[]{
                ".*CREDIT INTEREST.*",
                ".*SALARY.*"
        });
        return itemProcessor;
    }
}
