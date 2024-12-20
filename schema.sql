-- Create the Provider table
CREATE TABLE IF NOT EXISTS provider (
    id TEXT PRIMARY KEY, -- Unique identifier for the provider
    name TEXT NOT NULL   -- Provider's name
);

-- Create the Availability table
CREATE TABLE IF NOT EXISTS availability (
    id TEXT PRIMARY KEY,                          -- Unique identifier for availability
    provider_id TEXT NOT NULL,                    -- Foreign key to the provider
    start_time DATETIME NOT NULL,                 -- Start time of availability
    end_time DATETIME NOT NULL,                   -- End time of availability
    FOREIGN KEY (provider_id) REFERENCES provider (id), -- Enforce provider reference
    UNIQUE (provider_id, start_time, end_time)    -- Ensure no duplicate availability slots
);

-- Create the Slot table
CREATE TABLE IF NOT EXISTS slot (
    id TEXT PRIMARY KEY,                          -- Unique identifier for the slot
    availability_id TEXT NOT NULL,                -- Foreign key to the availability
    provider_id TEXT NOT NULL,                    -- Redundant provider reference for easier querying
    start_time DATETIME NOT NULL,                 -- Start time of the slot
    end_time DATETIME NOT NULL,                   -- End time of the slot
    status TEXT CHECK (status IN ('Available', 'Reserved', 'Confirmed')), -- Slot status
    FOREIGN KEY (availability_id) REFERENCES availability (id), -- Enforce availability reference
    FOREIGN KEY (provider_id) REFERENCES provider (id), -- Ensure slot references provider
    UNIQUE (availability_id, start_time, end_time) -- Ensure no duplicate slots
);

-- Create the Reservation table
CREATE TABLE IF NOT EXISTS reservation (
    id TEXT PRIMARY KEY,                          -- Unique identifier for the reservation
    client_id TEXT NOT NULL,                      -- Identifier for the client making the reservation
    provider_id TEXT NOT NULL,                    -- Foreign key to the provider
    availability_id TEXT NOT NULL,                -- Foreign key to the availability
    start_time DATETIME NOT NULL,                 -- Start time of the reservation
    end_time DATETIME NOT NULL,                   -- End time of the reservation
    reservation_expiry DATETIME,                  -- Expiry time for reservations
    status TEXT CHECK (status IN ('Pending', 'Confirmed')), -- Reservation status
    FOREIGN KEY (provider_id) REFERENCES provider (id), -- Enforce provider reference
    FOREIGN KEY (availability_id) REFERENCES availability (id) -- Enforce availability reference
);

-- Create Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_provider_id_availability ON availability (provider_id);
CREATE INDEX IF NOT EXISTS idx_availability_id_slot ON slot (availability_id);
CREATE INDEX IF NOT EXISTS idx_slot_provider_id ON slot (provider_id);
CREATE INDEX IF NOT EXISTS idx_client_id_reservation ON reservation (client_id);
CREATE INDEX IF NOT EXISTS idx_provider_id_reservation ON reservation (provider_id);
CREATE INDEX IF NOT EXISTS idx_availability_id_reservation ON reservation (availability_id);

-- Add a foreign key constraint to ensure availability references an existing provider
ALTER TABLE availability
ADD CONSTRAINT fk_provider
FOREIGN KEY (provider_id) REFERENCES provider(id)
ON DELETE CASCADE;