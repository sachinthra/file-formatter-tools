import { imageState, initialState } from '../hooks/useImageProcessor';
import type { JSX } from 'preact';

export function ImageUploader() {
  const handleFileChange = (event: JSX.TargetedEvent<HTMLInputElement, Event>) => {
    const target = event.currentTarget;
    if (target.files) {
      processFile(target.files[0]);
    }
  };

  const handleDragOver = (event: JSX.TargetedDragEvent<HTMLDivElement>) => {
    event.preventDefault();
    imageState.value = { ...imageState.value, isDragging: true };
  };

  const handleDragLeave = () => {
    imageState.value = { ...imageState.value, isDragging: false };
  };

  const handleDrop = (event: JSX.TargetedDragEvent<HTMLDivElement>) => {
    event.preventDefault();
    imageState.value = { ...imageState.value, isDragging: false };

    const file = event.dataTransfer?.files[0];
    if (file) {
      processFile(file);
    }
  };

  const processFile = (file: File) => {
    const originalSize = (file.size / 1024).toFixed(2);
    const img = new Image();

    imageState.value = { ...imageState.value, imageFile: file };

    img.onload = () => {
      imageState.value = {
        ...imageState.value,
        originalWidth: img.width,
        originalHeight: img.height,
        originalSize,
        width: img.width.toString(),
        height: img.height.toString(),
        maxSize: originalSize.toString(),
      };
    };

    img.src = URL.createObjectURL(file);
  };

  return (
    <div
      class={`border-2 border-dashed rounded-md p-8 text-center mb-4 cursor-pointer transition-all duration-200
        bg-blue-50/5 hover:border-blue-700 hover:bg-blue-50/10
        ${imageState.value.isDragging ? 'border-green-600 bg-green-100/10' : 'border-blue-500'}
        focus:outline-none`}
      role="button"
      tabIndex={0}
      onChange={handleFileChange}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      onClick={() => document.getElementById('image')?.click()}
    >
      {imageState.value.imageFile ? (
        <div class="flex items-center justify-center gap-2 text-gray-200">
          <span>Selected File: {imageState.value.imageFile.name}</span>
          <button
            type="button"
            class="ml-2 text-red-400 hover:text-red-600 text-lg font-bold px-1 rounded-full focus:outline-none"
            aria-label="Remove image"
            onClick={e => {
              e.stopPropagation();
              imageState.value = { ...initialState };
            }}
          >
            X
          </button>
        </div>
      ) : (
        <p class="text-gray-400">Drag and drop an image here, or click to select one.</p>
      )}
    </div>
  );
}