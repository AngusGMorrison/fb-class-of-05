-- Create the ENUM types on which tables depend.

CREATE TYPE "looking_for" AS ENUM (
  'friendship',
  'dating',
  'relationship',
  'random_play',
  'whatever'
);

CREATE TYPE "political_views" AS ENUM (
  'very_conservative',
  'conservative',
  'moderate',
  'liberal',
  'very_liberal'
);

CREATE TYPE "interested_in" AS ENUM (
  'men',
  'women',
  'both'
);

CREATE TYPE "relationship_status" AS ENUM (
  'single',
  'relationship',
  'engaged',
  'married',
  'complicated',
  'open'
);

CREATE TYPE "interest_category" AS ENUM (
  'general',
  'clubs_and_jobs',
  'movies',
  'music',
  'books'
);

CREATE TYPE "friend_request_status" AS ENUM (
  'accepted',
  'declined',
  'pending'
);
