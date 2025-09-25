-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create rule_types table
CREATE TABLE IF NOT EXISTS rule_types (
    id BIGSERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create rules table with vector embedding column
CREATE TABLE IF NOT EXISTS rules (
    id BIGSERIAL PRIMARY KEY,
    rule_type_id BIGINT NOT NULL REFERENCES rule_types(id) ON DELETE CASCADE,
    content JSONB NOT NULL,
    embedding vector(1536), -- OpenAI ada-002 embedding size, adjust as needed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_rules_rule_type_id ON rules(rule_type_id);
CREATE INDEX IF NOT EXISTS idx_rules_embedding ON rules USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_rule_types_updated_at 
    BEFORE UPDATE ON rule_types 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rules_updated_at 
    BEFORE UPDATE ON rules 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert some sample rule types
INSERT INTO rule_types (name) VALUES 
    ('validation'),
    ('transformation'),
    ('filtering'),
    ('business_logic')
ON CONFLICT (name) DO NOTHING;

-- Sample rules (optional - for testing)
INSERT INTO rules (rule_type_id, content) VALUES 
    (1, '{"description": "Validate email format", "pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$", "required": true}'),
    (2, '{"description": "Transform phone number", "format": "remove_spaces_and_dashes", "normalize": "+1"}'),
    (3, '{"description": "Filter adult content", "blacklist": ["adult", "explicit"], "threshold": 0.8}'),
    (4, '{"description": "Calculate discount", "conditions": {"min_amount": 100}, "discount_percent": 10}')
ON CONFLICT DO NOTHING;