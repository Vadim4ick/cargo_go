CREATE TABLE cargos (
  id SERIAL PRIMARY KEY,

  cargoNumber         TEXT UNIQUE NOT NULL,
  date                 TIMESTAMP,
  loadUnloadDate     TIMESTAMP,
  driver               TEXT NOT NULL,
  transportationInfo  TEXT NOT NULL,
  payoutAmount        NUMERIC,
  payoutDate          TIMESTAMP,
  paymentStatus       TEXT,
  payoutTerms         TEXT,

  "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  truckId INTEGER NOT NULL,
  CONSTRAINT fk_truck FOREIGN KEY (truckId) REFERENCES trucks(id) ON DELETE CASCADE
);
