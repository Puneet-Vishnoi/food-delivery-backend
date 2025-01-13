CREATE TABLE payment_requests (
    order_id VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,,
    amount NUMERIC(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL
);

CREATE TABLE payment_responses (
    id VARCHAR(255)NOT NULL PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    razorpay_key VARCHAR(255) NOT NULL
);

CREATE TABLE verify_payment_requests (
    order_id VARCHAR(255) NOT NULL,
    payment_id VARCHAR(255) NOT NULL,
    signature VARCHAR(255) NOT NULL,
    PRIMARY KEY (order_id, payment_id)
    FOREIGN KEY (order_id) REFERENCES payment_requests(order_id) ON DELETE CASCADE, 
    FOREIGN KEY (payment_id) REFERENCES payment_responses(id) ON DELETE CASCADE
);
