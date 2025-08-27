// MongoDB initialization script
// This script runs when MongoDB container starts for the first time

// Switch to pets_search database
db = db.getSiblingDB('pets_search');

// Create collections with basic schema validation
db.createCollection('users', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['email', 'created_at'],
            properties: {
                email: {
                    bsonType: 'string',
                    description: 'User email - required'
                },
                phone: {
                    bsonType: 'string',
                    description: 'User phone number'
                },
                name: {
                    bsonType: 'string',
                    description: 'User display name'
                },
                created_at: {
                    bsonType: 'date',
                    description: 'Creation timestamp - required'
                },
                updated_at: {
                    bsonType: 'date',
                    description: 'Last update timestamp'
                }
            }
        }
    }
});

db.createCollection('listings', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['user_id', 'type', 'title', 'status', 'created_at'],
            properties: {
                user_id: {
                    bsonType: 'objectId',
                    description: 'Reference to user - required'
                },
                type: {
                    bsonType: 'string',
                    enum: ['lost', 'found', 'adopt'],
                    description: 'Listing type - required'
                },
                title: {
                    bsonType: 'string',
                    description: 'Listing title - required'
                },
                description: {
                    bsonType: 'string',
                    description: 'Listing description'
                },
                city: {
                    bsonType: 'string',
                    description: 'City where pet was lost/found'
                },
                location: {
                    bsonType: 'string',
                    description: 'Specific location details'
                },
                contact_phone: {
                    bsonType: 'string',
                    description: 'Contact phone number'
                },
                contact_tg: {
                    bsonType: 'string',
                    description: 'Telegram contact'
                },
                status: {
                    bsonType: 'string',
                    enum: ['draft', 'active', 'archived'],
                    description: 'Listing status - required'
                },
                slug: {
                    bsonType: 'string',
                    description: 'URL slug for public page'
                },
                images: {
                    bsonType: 'array',
                    description: 'Array of image URLs'
                },
                created_at: {
                    bsonType: 'date',
                    description: 'Creation timestamp - required'
                },
                updated_at: {
                    bsonType: 'date',
                    description: 'Last update timestamp'
                }
            }
        }
    }
});

db.createCollection('events', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['listing_id', 'type', 'created_at'],
            properties: {
                user_id: {
                    bsonType: 'objectId',
                    description: 'User who triggered the event (optional)'
                },
                listing_id: {
                    bsonType: 'objectId',
                    description: 'Related listing - required'
                },
                type: {
                    bsonType: 'string',
                    enum: ['view', 'qr_scan', 'contact_click', 'phone_click'],
                    description: 'Event type - required'
                },
                payload: {
                    bsonType: 'object',
                    description: 'Additional event data'
                },
                ip_address: {
                    bsonType: 'string',
                    description: 'Client IP address'
                },
                user_agent: {
                    bsonType: 'string',
                    description: 'Client user agent'
                },
                created_at: {
                    bsonType: 'date',
                    description: 'Event timestamp - required'
                }
            }
        }
    }
});

// Create indexes for better performance
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ phone: 1 });

db.listings.createIndex({ user_id: 1 });
db.listings.createIndex({ type: 1 });
db.listings.createIndex({ status: 1 });
db.listings.createIndex({ slug: 1 }, { unique: true, sparse: true });
db.listings.createIndex({ city: 1 });
db.listings.createIndex({ created_at: -1 });

db.events.createIndex({ listing_id: 1 });
db.events.createIndex({ type: 1 });
db.events.createIndex({ created_at: -1 });
db.events.createIndex({ user_id: 1 });

print('MongoDB initialization completed successfully!');
