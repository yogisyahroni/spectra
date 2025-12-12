# SPECTRA: INTELLIGENT FIBER INFRASTRUCTURE SYSTEM

**Document Type:** Master Technical Blueprint (End-to-End)  
**Version:** 3.1 (Refined - Enhanced Structure and Completeness)  
**Objective:** To provide an integrated system encompassing Geographic Information System (GIS), Network Management System (NMS), splicing management, and field operations for efficient fiber optic infrastructure management. This blueprint ensures seamless visualization, real-time monitoring, and operational efficiency.

---

## 1. INTRODUCTION AND SCOPE

### A. Overview
SPECTRA is a comprehensive platform designed to manage fiber optic networks by integrating spatial data visualization, inventory tracking, network monitoring, and field technician workflows. It leverages advanced mapping technologies to render infrastructure assets dynamically, supports real-time status updates from NMS, and facilitates precise splicing and connection management.

### B. Scope
- **In Scope:** GIS-based asset mapping, NMS integration for status monitoring, splicing and connection logic, field mobile applications for technicians, inventory management, and basic security measures.
- **Out of Scope:** Advanced analytics (e.g., predictive maintenance using AI), integration with third-party billing systems, or hardware procurement specifications.
- **Target Users:** Network Operations Center (NOC) administrators, field technicians, and infrastructure planners.

### C. Assumptions and Dependencies
- Access to accurate GPS data and NMS APIs from vendors like Huawei and ZTE.
- Stable internet connectivity for web and mobile components, with offline fallback for field apps.
- Compliance with data privacy regulations (e.g., GDPR or equivalent local standards).

---

## 2. TECHNOLOGY STACK ARCHITECTURE
*Selected based on performance for heavy visual rendering and concurrent NMS data handling.*

| Layer | Component | Technology Selection | Specification & Justification |
|-------|-----------|-----------------------|-------------------------------|
| **GIS Engine** | **Map Renderer** | **Mapbox GL JS (Web) & Mapbox SDK (Mobile)** | Utilizes WebGL for rendering vector icons and lines. Significantly lighter than Google Maps API when handling over 50,000 assets. Supports data-driven styling, enabling automatic color changes based on database values. |
| **Database** | **Spatial DB** | **PostgreSQL + PostGIS** | Stores geometric data (e.g., POINT, LINESTRING) and performs spatial queries (e.g., radius-based searches for nearby assets). |
| **Backend** | **Core API** | **Golang (Go)** | Provides high-performance concurrency to manage thousands of NMS requests and real-time technician position updates. |
| **Frontend** | **Web Admin** | **React.js + Vite** | Modern Single Page Application (SPA) framework for responsive and efficient user interfaces. |
| **Mobile** | **Field App** | **Flutter** | Cross-platform development with native map rendering performance (60 FPS) and offline map support. |
| **Messaging** | **Broker** | **Redis + RabbitMQ** | Handles caching for real-time status updates, preventing continuous map reloads and ensuring low-latency operations. |

---

## 3. MAP VISUALIZATION SPECIFICATION (EXPLICIT)
*This section details the translation of coordinate data into visual elements on the Mapbox map, ensuring clarity and usability.*

### A. Symbology (Asset Shapes and Icons)
Each asset type has a unique visual representation for easy identification.

| Asset Type | Visual Shape (Icon) | Default Color | Visual Description |
|------------|---------------------|---------------|--------------------|
| **OLT** (Central Hub) | 🏢 **Large Square Building** | 🔵 **Dark Blue** | Icon size is 2x larger than others; name label always visible for quick reference. |
| **ODC** (Cabinet) | 📦 **Square Box** | 🟢 **Dark Green** | Includes a small badge displaying capacity (e.g., "288") for at-a-glance assessment. |
| **ODP** (Pole) | ⚪ **Circle / Dot** | 🔵 **Light Blue** | Solid circle with white border; size scales down during zoom-out to maintain map clarity. |
| **Closure** (Joint) | 💊 **Capsule / Oval** | 🟠 **Orange** | Positioned midway along cable paths, not at endpoints like customer locations. |
| **Customer** (Home) | 🏠 **House Icon** | ⚫ **Black (Outline)** | Small icon; fill color changes based on status (Green=Online, Red=Loss of Signal). |
| **Slack** (Coil) | ➰ **Loop Icon** | 🟤 **Brown** | Indicates excess cable coils at poles for maintenance planning. |

### B. Cable Visualization (Cable Paths)
Cables are rendered using GeoJSON LineString features, with thickness and style indicating type.

- **Backbone Cable (Feeder):**
  - *Color:* 🔴 **Red**
  - *Style:* Solid thick line (Width: 4px).
  - *Visibility:* Prominent at low zoom levels (city-wide views).
- **Distribution Cable:**
  - *Color:* 🔵 **Blue**
  - *Style:* Solid medium line (Width: 2px).
- **Drop Cable (To Home):**
  - *Color:* ⚫ **Black**
  - *Style:* Dashed line.
  - *Visibility:* Appears only at maximum zoom-in (Level 16+).

### C. Status Indication (Real-Time Feedback)
The map dynamically updates visuals based on NMS data for proactive monitoring.

1. **Loss of Signal (LOS) Status:**
   - *Effect:* "Pulsing Halo" around icons (expanding/contracting transparent red circle).
   - *Purpose:* Immediately draws NOC operator attention to issues.
2. **Full Capacity Status:**
   - *Effect:* ODP icon changes to solid red.
   - *Trigger:* When used_ports equals total_ports.

### D. Zoom Level Strategy (Clustering)
To prevent visual clutter, display logic adapts to zoom levels:

- **Zoom Level 10-12 (City View):** Display only OLTs and red backbone lines; hide ODPs.
- **Zoom Level 13-15 (District View):** Show ODCs and blue distribution cables; cluster ODPs (e.g., circle with "50 ODP").
- **Zoom Level 16+ (Street View):** Reveal all elements (ODPs, closures, poles, homes, drop cables); clusters expand into individual points.

---

## 4. PRODUCT REQUIREMENTS DOCUMENT (PRD)

### A. Core Features

#### Module 1: Inventory Management
- **Physical Data Input:** User forms for entering coordinates (Latitude/Longitude), asset types, and cable specifications, with validation for data integrity.
- **Splicing Logic:** Interactive UI matrix for connecting cores (e.g., straight, cross, or branch splicing), including loss calculation and validation.
- **Trace Route Logic:** Backend algorithm to trace physical paths from customer homes to OLT, generating reports and visualizations.

#### Module 2: NMS Integration
- **Synchronization:** Go-based daemon for periodic polling of ONT status data from Huawei/ZTE systems.
- **Mapping:** Automated matching of NMS serial numbers to GIS customer IDs for unified data views.
- **Alarm Handling:** Visual notifications on the map for poor Rx Power (threshold: > -27 dBm), with email/SMS alerts.

#### Module 3: Field App (Technician Tools)
- **Radius Lock:** "Check-In" button disabled if device GPS is beyond 50 meters from the asset location.
- **Navigation:** Integration with Google Maps for directions upon selecting an asset.
- **Offline First Design:** Local storage of data during network loss, with automatic synchronization upon reconnection.

### B. Non-Functional Requirements
- **Performance:** System must handle 1,000 concurrent users with <2-second map refresh times.
- **Scalability:** Horizontal scaling via containerization (e.g., Docker/Kubernetes).
- **Usability:** Intuitive interfaces with accessibility features (e.g., high-contrast modes).
- **Reliability:** 99.9% uptime, with automated backups and failover mechanisms.

### C. Prioritization and Milestones
- **MVP (Minimum Viable Product):** Core GIS mapping, inventory input, and basic NMS sync (Q1 2026).
- **Full Release:** All modules, including field app and advanced visualizations (Q3 2026).

---

## 5. ENTITY RELATIONSHIP DIAGRAM (ERD)
The following SQL schema defines the core database structure, emphasizing spatial and relational integrity.

```sql
-- Nodes (Map Points)
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,  -- e.g., "ODP-JATI-01"
    type VARCHAR(20) CHECK (type IN ('OLT', 'ODC', 'ODP', 'CLOSURE', 'POLE', 'CUSTOMER')),
    coordinate GEOMETRY(POINT, 4326) NOT NULL,  -- Mandatory PostGIS column
    address TEXT,
    capacity_ports INT DEFAULT 8,
    model VARCHAR(50),
    status VARCHAR(20) DEFAULT 'ACTIVE'  -- ACTIVE, MAINTENANCE, PLAN
);

-- Cables
CREATE TABLE cables (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),  -- e.g., "Feeder 01"
    type VARCHAR(20) CHECK (type IN ('ADSS', 'DUCT', 'DROP')),
    core_count INT NOT NULL,  -- e.g., 12, 24, 48, 96, 144
    length_meter FLOAT,
    origin_node_id BIGINT REFERENCES nodes(id),
    dest_node_id BIGINT REFERENCES nodes(id),
    path_geometry GEOMETRY(LINESTRING, 4326),  -- For line visualization
    color_hex VARCHAR(7) DEFAULT '#000000'  -- Custom color override
);

-- Cable Cores (Internal Details)
CREATE TABLE cable_cores (
    id BIGSERIAL PRIMARY KEY,
    cable_id BIGINT REFERENCES cables(id) ON DELETE CASCADE,
    core_index INT NOT NULL,  -- 1-144
    tube_color VARCHAR(30),  -- e.g., "Blue"
    core_color VARCHAR(30),  -- e.g., "Blue"
    status VARCHAR(20) DEFAULT 'VACANT'
);

-- Connections (Splicing Links)
CREATE TABLE connections (
    id BIGSERIAL PRIMARY KEY,
    location_node_id BIGINT REFERENCES nodes(id),  -- Closure location
    input_type VARCHAR(10),  -- CORE / PORT
    input_id BIGINT,
    output_type VARCHAR(10),  -- CORE / PORT
    output_id BIGINT,
    loss_db FLOAT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Customers
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT REFERENCES nodes(id),  -- GIS link
    name VARCHAR(100),
    ont_sn VARCHAR(50) UNIQUE,  -- NMS link
    current_status VARCHAR(20),  -- ONLINE, LOS, POWER_OFF
    last_rx_power FLOAT
);

6. UI/UX FLOW (USER JOURNEYS)
A. Scenario: Technician New Installation (Mobile App)

Open Map: Technician views their position (blue dot) and nearby ODPs (blue circles).
Select ODP: Tap nearest ODP.
Pop-Up Info: Displays "ODP-JATI-01 (Distance: 10m). Capacity: 5/8 Used."
View Ports: Click "View Ports" to show an 8-slot grid; slots 1-5 red (occupied), 6-8 green (available).
Action: Select slot 6 → Activate camera → Scan modem barcode → Capture home photo → Submit.

B. Scenario: NOC Monitoring (Web Dashboard)

Global View: Map shows red backbone lines; no alarms present.
Event Trigger: Pulsing red circle appears in South Jakarta area.
Investigation: Zoom into red area; ODC flashes red, with faded red cable paths below.
Info Box: "Mass Outage Detected. 45 Users LOS. Suspected: Feeder Cable Cut."
Action: Right-click ODC → "Create Trouble Ticket" → Assign to South Jakarta maintenance team.


7. SECURITY AND DEPLOYMENT
A. Security Measures

API Key Protection: Restrict Mapbox tokens to application domains/bundle IDs to prevent unauthorized usage.
Data Encryption: Apply AES-256 encryption to sensitive database fields (e.g., OLT passwords, customer data).
Access Control: Host web admin in a private subnet, requiring VPN for access; implement role-based access control (RBAC).
Auditing: Log all actions with timestamps and user IDs for compliance and forensic analysis.

B. Deployment Strategy

Environment: Cloud-based (e.g., AWS or Azure) with auto-scaling groups.
CI/CD Pipeline: Use GitHub Actions or Jenkins for automated testing and deployment.
Monitoring: Integrate tools like Prometheus and Grafana for performance metrics and alerts.

This refined blueprint enhances the original structure by adding an introduction, scope, non-functional requirements, and deployment details, ensuring a comprehensive and professional document.