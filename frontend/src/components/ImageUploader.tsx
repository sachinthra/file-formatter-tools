import { imageState } from '../hooks/useImageProcessor';
import type { JSX } from 'preact';
import styles from './ImageUploader.module.css';

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
      class={`${styles.dropZone} ${imageState.value.isDragging ? styles.isDragging : ''}`}
      role="button"
      tabIndex={0}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      {imageState.value.imageFile ? (
        <p>Selected File: {imageState.value.imageFile.name}</p>
      ) : (
        <p>Drag and drop an image here, or click to select one.</p>
      )}
      <input
        type="file"
        id="image"
        accept="image/*"
        onChange={handleFileChange}
        style={{ display: 'none' }}
      />
    </div>
  );
}