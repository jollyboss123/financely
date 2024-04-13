package org.jolly.financely;

import org.jolly.financely.money.Money;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.math.RoundingMode;
import java.util.Currency;
import java.util.Locale;

@SpringBootApplication
public class FinancelyApplication {

    public static void main(String[] args) {
        SpringApplication.run(FinancelyApplication.class, args);
        Money.init(Currency.getInstance(Locale.of("en", "MY")), RoundingMode.HALF_EVEN);
    }

}
