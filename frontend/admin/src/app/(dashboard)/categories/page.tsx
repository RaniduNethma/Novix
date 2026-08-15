'use client';

import { useEffect, useState } from 'react';
import { Category } from '@/types';
import { categoryService } from '@/services/category.service';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Spinner } from '@/components/ui/Spinner';
import { Plus, Trash2 } from 'lucide-react';
import toast from 'react-hot-toast';

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form, setForm] = useState({
    name: '',
    description: '',
    slug: '',
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const loadCategories = async () => {
    setIsLoading(true);
    try {
      const data = await categoryService.getCategories(0, 50);
      setCategories(data.content);
    } catch {
      toast.error('Failed to load categories');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadCategories();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      await categoryService.createCategory(form.name, form.description, form.slug);
      toast.success('Category created!');
      setIsModalOpen(false);
      setForm({ name: '', description: '', slug: '' });
      loadCategories();
    } catch {
      toast.error('Failed to create category');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (id: number, name: string) => {
    if (!confirm(`Delete category "${name}"?`)) return;
    try {
      await categoryService.deleteCategory(id);
      toast.success('Category deleted');
      loadCategories();
    } catch {
      toast.error('Failed to delete category');
    }
  };

  const autoSlug = (name: string) =>
    name
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^a-z0-9-]/g, '');

  return (
    <div className="space-y-6">
      <div className="flex justify-end">
        <Button variant="primary" onClick={() => setIsModalOpen(true)}>
          <Plus className="h-4 w-4 mr-1" />
          New Category
        </Button>
      </div>

      <Card className="p-0 overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center py-16">
            <Spinner className="h-8 w-8" />
          </div>
        ) : (
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-800">
                <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">Name</th>
                <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">Slug</th>
                <th className="text-left text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                  Description
                </th>
                <th className="text-right text-xs font-medium text-gray-400 px-6 py-3 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-800">
              {categories.map((cat) => (
                <tr key={cat.id} className="hover:bg-gray-800/50">
                  <td className="px-6 py-4 text-sm font-medium text-white">{cat.name}</td>
                  <td className="px-6 py-4 text-sm text-gray-400 font-mono">{cat.slug}</td>
                  <td className="px-6 py-4 text-sm text-gray-400">{cat.description || '—'}</td>
                  <td className="px-6 py-4 text-right">
                    <Button size="sm" variant="danger" onClick={() => handleDelete(cat.id, cat.name)}>
                      <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>

      {/* Create Category Modal */}
      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Create Category">
        <form onSubmit={handleCreate} className="space-y-4">
          <Input
            label="Name"
            placeholder="Entertainment"
            value={form.name}
            onChange={(e) =>
              setForm({
                ...form,
                name: e.target.value,
                slug: autoSlug(e.target.value),
              })
            }
            required
          />
          <Input
            label="Slug"
            placeholder="entertainment"
            value={form.slug}
            onChange={(e) =>
              setForm({
                ...form,
                slug: e.target.value,
              })
            }
            required
          />
          <Input
            label="Description"
            placeholder="Optional description"
            value={form.description}
            onChange={(e) =>
              setForm({
                ...form,
                description: e.target.value,
              })
            }
          />
          <div className="flex gap-3 pt-2">
            <Button type="button" variant="secondary" className="flex-1" onClick={() => setIsModalOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" variant="primary" className="flex-1" isLoading={isSubmitting}>
              Create
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
