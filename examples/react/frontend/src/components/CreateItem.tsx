import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api, APIError } from '../api/client';
import { CreateItemRequest } from '../types/api';

export const CreateItem: React.FC = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState<CreateItemRequest>({
    name: '',
    description: '',
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      await api.createItem(formData);
      navigate('/');
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError('Failed to create item');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="container">
      <div className="header">
        <h1>Create New Item</h1>
        <Link to="/" className="button">
          Back to List
        </Link>
      </div>

      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit} className="form">
        <div className="form-group">
          <label htmlFor="name">Name</label>
          <input
            id="name"
            type="text"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Enter item name"
            required
            disabled={submitting}
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description</label>
          <textarea
            id="description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Enter item description"
            rows={4}
            required
            disabled={submitting}
          />
        </div>

        <div className="form-actions">
          <button type="submit" className="button button-primary" disabled={submitting}>
            {submitting ? 'Creating...' : 'Create Item'}
          </button>
          <Link to="/" className="button">
            Cancel
          </Link>
        </div>
      </form>
    </div>
  );
};

