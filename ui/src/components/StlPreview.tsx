import { useMemo } from "react";
import { Canvas } from "@react-three/fiber";
import { OrbitControls } from "@react-three/drei";
import { STLLoader } from "three-stdlib";
import * as THREE from "three";

type STLViewerProps = {
  stlData: string;
  color: string;
};

const STLModel = ({ stlData, color }: { stlData: string; color: string }) => {
  const geometry = useMemo(() => {
    const loader = new STLLoader();
    // Convert base64 to ArrayBuffer
    const binary = atob(stlData);
    const arrayBuffer = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      arrayBuffer[i] = binary.charCodeAt(i);
    }
    const geom = loader.parse(arrayBuffer.buffer);
    geom.computeVertexNormals();
    return geom;
  }, [stlData]);

  const scale = useMemo(() => {
    geometry.computeBoundingBox();
    const box = geometry.boundingBox;
    if (!box) return 1;
    const size = new THREE.Vector3();
    box.getSize(size);
    const maxDim = Math.max(size.x, size.y, size.z);
    return maxDim > 0 ? 1 / maxDim : 1;
  }, [geometry]);

  const material = useMemo(() => {
    return new THREE.MeshStandardMaterial({ color });
  }, [color]);

  return (
    <mesh geometry={geometry} scale={[scale, scale, scale]}>
      <primitive object={material} attach="material" />
    </mesh>
  );
};

export const STLPreview = ({ stlData, color }: STLViewerProps) => {
  return (
    <Canvas style={{ height: 400, width: "100%" }}>
      <ambientLight intensity={0.5} />
      <directionalLight position={[2, 2, 2]} intensity={0.8} />
      <axesHelper args={[2]} />
      <STLModel stlData={stlData} color={color} />
      <OrbitControls
        enableDamping
        dampingFactor={0.05}
        rotateSpeed={0.5}
        maxDistance={10}
        minDistance={0.5}
      />
    </Canvas>
  );
};