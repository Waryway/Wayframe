import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api, APIError } from '../api/client';
import { Item } from '../types/api';

export const ItemList: React.FC = () => {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const perPage = 10;

  useEffect(() => {
    loadItems();
  }, [page]);

  const loadItems = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.getItems(page, perPage);
      setItems(response.data);
      setTotal(response.total);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError('Failed to load items');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this item?')) {
      return;
    }

    try {
      await api.deleteItem(id);
      loadItems();
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
        <div className="loading">Loading items...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <div className="error">Error: {error}</div>
        <button onClick={loadItems}>Retry</button>
      </div>
    );
  }

  const totalPages = Math.ceil(total / perPage);

  return (
    <div className="container">
      <div className="header">
        <h1>Items</h1>
        <Link to="/create" className="button button-primary">
          Create New Item
        </Link>
      </div>

      {items.length === 0 ? (
        <div className="empty-state">
          <p>No items found. Create your first item!</p>
        </div>
      ) : (
        <>
          <div className="items-grid">
            {items.map((item) => (
              <div key={item.id} className="item-card">
                <h3>{item.name}</h3>
                <p>{item.description}</p>
                <div className="item-meta">
                  <small>Created: {new Date(item.created_at).toLocaleDateString()}</small>
                </div>
                <div className="item-actions">
                  <Link to={`/items/${item.id}`} className="button button-sm">
                    View
                  </Link>
                  <button
                    onClick={() => handleDelete(item.id)}
                    className="button button-sm button-danger"
                  >
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="pagination">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="button button-sm"
              >
                Previous
              </button>
              <span>
                Page {page} of {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="button button-sm"
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

