package org.jolly.financely;

/**
 * @author jolly
 */
public interface Auditable {
    Audit getAudit();
    void setAudit(Audit audit);
}
