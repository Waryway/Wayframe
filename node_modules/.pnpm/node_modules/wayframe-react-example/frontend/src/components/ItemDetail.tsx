import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { api, APIError } from '../api/client';
import { Item, UpdateItemRequest } from '../types/api';

export const ItemDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [item, setItem] = useState<Item | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [formData, setFormData] = useState<UpdateItemRequest>({});

  useEffect(() => {
    loadItem();
  }, [id]);

  const loadItem = async () => {
    if (!id) return;

    try {
      setLoading(true);
      setError(null);
      const data = await api.getItem(parseInt(id));
      setItem(data);
      setFormData({ name: data.name, description: data.description });
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError('Failed to load item');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;

    try {
      const updated = await api.updateItem(parseInt(id), formData);
      setItem(updated);
      setEditing(false);
    } catch (err) {
      if (err instanceof APIError) {
        alert(`Failed to update: ${err.message}`);
      } else {
        alert('Failed to update item');
      }
    }
  };

  const handleDelete = async () => {
    if (!id || !window.confirm('Are you sure you want to delete this item?')) {
      return;
    }

    try {
      await api.deleteItem(parseInt(id));
      navigate('/');
    } catch (err) {
      if (err instanceof APIError) {
        alert(`Failed to delete: ${err.message}`);
      } else {
        alert('Failed to delete item');
      }
    }
  };

  if (loading) {
    return (
      <div className="container">
        <div className="loading">Loading item...</div>
      </div>
    );
  }

  if (error || !item) {
    return (
      <div className="container">
        <div className="error">Error: {error || 'Item not found'}</div>
        <Link to="/" className="button">
          Back to List
        </Link>
      </div>
    );
  }

  return (
    <div className="container">
      <div className="header">
        <h1>Item Details</h1>
        <Link to="/" className="button">
          Back to List
        </Link>
      </div>

      {editing ? (
        <form onSubmit={handleUpdate} className="form">
          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input
              id="name"
              type="text"
              value={formData.name || ''}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              value={formData.description || ''}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows={4}
              required
            />
          </div>

          <div className="form-actions">
            <button type="submit" className="button button-primary">
              Save Changes
            </button>
            <button
              type="button"
              onClick={() => {
                setEditing(false);
                setFormData({ name: item.name, description: item.description });
              }}
              className="button"
            >
              Cancel
            </button>
          </div>
        </form>
      ) : (
        <div className="item-detail">
          <div className="detail-section">
            <h2>{item.name}</h2>
            <p>{item.description}</p>
          </div>

          <div className="detail-meta">
            <div className="meta-item">
              <strong>ID:</strong> {item.id}
            </div>
            <div className="meta-item">
              <strong>Created:</strong> {new Date(item.created_at).toLocaleString()}
            </div>
            <div className="meta-item">
              <strong>Updated:</strong> {new Date(item.updated_at).toLocaleString()}
            </div>
          </div>

          <div className="detail-actions">
            <button onClick={() => setEditing(true)} className="button button-primary">
              Edit
            </button>
            <button onClick={handleDelete} className="button button-danger">
              Delete
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

