package com.thoughtworks.gauge.test.implementation;

import com.thoughtworks.gauge.Step;
import com.thoughtworks.gauge.Table;
import com.thoughtworks.gauge.TableRow;
import com.thoughtworks.gauge.test.common.builders.SpecificationBuilder;

public class Specification {

    @Step("Create specs <table>")
    public void createSpecs(Table specs) throws Exception {
        for (int i = 0; i < specs.getTableRows().size(); i++) {
            TableRow row = specs.getTableRows().get(i);
            createSpec(row.getCell("heading"), row.getCell("path"), "spec" + i);
        }
    }

    public void createSpec(String specName, String subFolder, String filename) throws Exception {
        new SpecificationBuilder().withScenarioName("A scenario").withSpecName(specName).withSteps(Steps.example())
                .withSubDirPath(subFolder).withFilename(filename).buildAndAddToProject(false);
    }

}
