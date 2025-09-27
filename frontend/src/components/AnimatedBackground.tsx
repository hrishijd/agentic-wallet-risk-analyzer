import { useEffect, useRef } from 'react';

interface Node {
  x: number;
  y: number;
  vx: number;
  vy: number;
  size: number;
  color: string;
  name: string;
}

interface Connection {
  from: number;
  to: number;
  opacity: number;
  pulse: number;
}

const AnimatedBackground = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationFrameRef = useRef<number>();
  const nodesRef = useRef<Node[]>([]);
  const connectionsRef = useRef<Connection[]>([]);

  // Blockchain colors
  const chainColors = [
    'hsl(205, 100%, 60%)', // Ethereum Blue
    'hsl(200, 100%, 65%)', // Light Blue
    'hsl(210, 100%, 70%)', // Lighter Blue
    'hsl(195, 100%, 75%)', // Cyan
    'hsl(215, 100%, 55%)', // Deep Blue
    'hsl(190, 100%, 60%)', // Aqua
  ];

  const chainNames = ['ETH', 'BSC', 'POLY', 'AVAX', 'SOL', 'ARB'];

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas size
    const resizeCanvas = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };
    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);

    // Initialize nodes
    const initNodes = () => {
      nodesRef.current = chainNames.map((name, i) => ({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        vx: (Math.random() - 0.5) * 0.5,
        vy: (Math.random() - 0.5) * 0.5,
        size: 8 + Math.random() * 4,
        color: chainColors[i],
        name: name,
      }));

      // Create connections between some nodes
      connectionsRef.current = [];
      for (let i = 0; i < nodesRef.current.length; i++) {
        for (let j = i + 1; j < nodesRef.current.length; j++) {
          if (Math.random() > 0.6) {
            connectionsRef.current.push({
              from: i,
              to: j,
              opacity: 0.3 + Math.random() * 0.4,
              pulse: Math.random() * Math.PI * 2,
            });
          }
        }
      }
    };

    initNodes();

    // Animation loop
    const animate = () => {
      ctx.fillStyle = 'rgba(215, 25, 8, 0.05)';
      ctx.fillRect(0, 0, canvas.width, canvas.height);

      const nodes = nodesRef.current;
      const connections = connectionsRef.current;

      // Update and draw connections
      connections.forEach((conn) => {
        const nodeA = nodes[conn.from];
        const nodeB = nodes[conn.to];
        
        conn.pulse += 0.05;
        const pulseOpacity = (Math.sin(conn.pulse) + 1) * 0.5;
        const opacity = conn.opacity * pulseOpacity * 0.6;

        // Draw connection line
        ctx.strokeStyle = `rgba(73, 166, 230, ${opacity})`;
        ctx.lineWidth = 1.5;
        ctx.beginPath();
        ctx.moveTo(nodeA.x, nodeA.y);
        ctx.lineTo(nodeB.x, nodeB.y);
        ctx.stroke();

        // Draw moving particles along connection
        const distance = Math.sqrt(
          Math.pow(nodeB.x - nodeA.x, 2) + Math.pow(nodeB.y - nodeA.y, 2)
        );
        const numParticles = Math.floor(distance / 60);
        
        for (let i = 0; i < numParticles; i++) {
          const progress = (conn.pulse * 0.1 + i / numParticles) % 1;
          const x = nodeA.x + (nodeB.x - nodeA.x) * progress;
          const y = nodeA.y + (nodeB.y - nodeA.y) * progress;
          
          ctx.fillStyle = `rgba(73, 166, 230, ${opacity * 2})`;
          ctx.beginPath();
          ctx.arc(x, y, 1.5, 0, Math.PI * 2);
          ctx.fill();
        }
      });

      // Update and draw nodes
      nodes.forEach((node) => {
        // Update position
        node.x += node.vx;
        node.y += node.vy;

        // Bounce off walls
        if (node.x <= 0 || node.x >= canvas.width) node.vx *= -1;
        if (node.y <= 0 || node.y >= canvas.height) node.vy *= -1;

        // Keep nodes in bounds
        node.x = Math.max(0, Math.min(canvas.width, node.x));
        node.y = Math.max(0, Math.min(canvas.height, node.y));

        // Draw node glow
        const gradient = ctx.createRadialGradient(
          node.x, node.y, 0,
          node.x, node.y, node.size * 2
        );
        gradient.addColorStop(0, node.color);
        gradient.addColorStop(1, 'transparent');
        
        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size * 2, 0, Math.PI * 2);
        ctx.fill();

        // Draw node core
        ctx.fillStyle = node.color;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size, 0, Math.PI * 2);
        ctx.fill();

        // Draw node border
        ctx.strokeStyle = 'rgba(255, 255, 255, 0.8)';
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.arc(node.x, node.y, node.size, 0, Math.PI * 2);
        ctx.stroke();

        // Draw node label
        ctx.fillStyle = 'rgba(210, 40, 95, 0.9)';
        ctx.font = 'bold 10px monospace';
        ctx.textAlign = 'center';
        ctx.fillText(node.name, node.x, node.y - node.size - 8);
      });

      animationFrameRef.current = requestAnimationFrame(animate);
    };

    animate();

    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
      window.removeEventListener('resize', resizeCanvas);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      className="fixed inset-0 w-full h-full pointer-events-none z-0"
      style={{ background: 'transparent' }}
    />
  );
};

export default AnimatedBackground;
