package org.jolly.financely.model;

/**
 * @author jolly
 */
public interface Auditable {
    Audit getAudit();
    void setAudit(Audit audit);
}
