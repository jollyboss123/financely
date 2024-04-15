package org.jolly.financely.batch.job;

import org.jolly.financely.batch.processor.DefaultExpenseProcessor;
import org.jolly.financely.model.Expense;
import org.jolly.financely.model.Money;
import org.springframework.batch.core.Job;
import org.springframework.batch.core.Step;
import org.springframework.batch.core.job.builder.JobBuilder;
import org.springframework.batch.core.launch.support.RunIdIncrementer;
import org.springframework.batch.core.repository.JobRepository;
import org.springframework.batch.core.step.builder.StepBuilder;
import org.springframework.batch.item.ItemProcessor;
import org.springframework.batch.item.ItemReader;
import org.springframework.batch.item.ItemWriter;
import org.springframework.batch.item.file.FlatFileItemReader;
import org.springframework.batch.item.file.LineMapper;
import org.springframework.batch.item.file.MultiResourceItemReader;
import org.springframework.batch.item.file.mapping.BeanWrapperFieldSetMapper;
import org.springframework.batch.item.file.mapping.DefaultLineMapper;
import org.springframework.batch.item.file.mapping.FieldSetMapper;
import org.springframework.batch.item.file.transform.DelimitedLineTokenizer;
import org.springframework.batch.item.file.transform.FieldSet;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.Resource;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.validation.BindException;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

/**
 * @author jolly
 */
@Configuration
public class ExpenseBatchJob {
    @Value("file:${file.path.expense}")
    private Resource[] resources;
    @Value("${fields.name.expense:#{null}}")
    private String[] fieldNames;
    @Value("${date.format.expense:#{null}}")
    private String dateFormat;
    private static final String JOB_NAME = "Expense.ETL.Job";
    private static final String PROCESSOR_TASK_NAME = "Expense.ETL.Job.file.load";

    @Bean
    public Job expenseJob(JobRepository jobRepository,
                          PlatformTransactionManager transactionManager,
                          ItemReader<Expense> expenseItemsReader,
                          ItemProcessor<Expense, Expense> expenseItemProcessor,
                          ItemWriter<Expense> expenseDBWriter) {
        Step step = new StepBuilder(PROCESSOR_TASK_NAME, jobRepository)
                .<Expense, Expense>chunk(100, transactionManager)
                .reader(expenseItemsReader)
                .processor(expenseItemProcessor)
                .writer(expenseDBWriter)
                .build();

        return new JobBuilder(JOB_NAME, jobRepository)
                .incrementer(new RunIdIncrementer())
                .start(step)
                .build();
    }

    @Bean
    public MultiResourceItemReader<Expense> expenseItemsReader(FlatFileItemReader<Expense> expenseItemReader) {
        MultiResourceItemReader<Expense> reader = new MultiResourceItemReader<>();
        reader.setResources(resources);
        reader.setStrict(false);
        reader.setDelegate(expenseItemReader);
        return reader;
    }

    @Bean
    public FlatFileItemReader<Expense> expenseItemReader() {
        FlatFileItemReader<Expense> reader = new FlatFileItemReader<>();
        reader.setName("expense.csv.reader");
        reader.setLinesToSkip(1);
        reader.setLineMapper(expenseLineMapper());
        reader.setStrict(false);
        return reader;
    }

    @Bean
    public LineMapper<Expense> expenseLineMapper() {
        DefaultLineMapper<Expense> lineMapper = new DefaultLineMapper<>();

        DelimitedLineTokenizer tokenizer = new DelimitedLineTokenizer();
        tokenizer.setDelimiter(",");
        tokenizer.setStrict(false);
        tokenizer.setNames(fieldNames);

        FieldSetMapper<Expense> fieldSetMapper = fieldSet -> {
            DateTimeFormatter dateTimeFormatter = DateTimeFormatter.ofPattern(dateFormat);
            LocalDateTime date = LocalDateTime.parse(fieldSet.readString(0), dateTimeFormatter);
            Money price = Money.of(BigDecimal.valueOf(Double.parseDouble(fieldSet.readString(2))));

            return new Expense.Builder(date, price, fieldSet.readString(3))
                    .description(fieldSet.readString(1))
                    .build();
        };

        lineMapper.setLineTokenizer(tokenizer);
        lineMapper.setFieldSetMapper(fieldSetMapper);

        return lineMapper;
    }

    @Bean
    public DefaultExpenseProcessor expenseItemProcessor(@Qualifier("defaultExpenseProcessor") DefaultExpenseProcessor itemProcessor) {
        itemProcessor.setItemsToSkip(new String[]{
            ".*Income.*"
        });
        itemProcessor.setCreditTransfer(new String[]{

        });

        return itemProcessor;
    }
}
