package com.thoughtworks.gauge.test.common.builders;

import com.thoughtworks.gauge.Table;
import com.thoughtworks.gauge.test.common.Concept;

import static com.thoughtworks.gauge.test.common.GaugeProject.getCurrentProject;

public class ConceptBuilder {
    private String conceptName;
    private Table steps;
    private String subDirPath;

    public ConceptBuilder withName(String conceptName) {
        this.conceptName = conceptName;
        return this;
    }

    public ConceptBuilder withSteps(Table steps) {
        this.steps = steps;
        return this;
    }

    public ConceptBuilder withSubDirPath(String subDirPath) {
        this.subDirPath = subDirPath;
        return this;
    }

    public Concept build() throws Exception {
        Concept concept = getCurrentProject().createConcept(subDirPath, conceptName, steps);
        getCurrentProject().addConcepts(concept);
        return concept;
    }
}
