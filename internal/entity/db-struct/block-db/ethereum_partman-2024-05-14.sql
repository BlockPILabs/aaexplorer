
SELECT partman.create_parent(p_parent_table => 'public.block_data_decode',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'weekly',
                             p_premake => 5,
                             p_start_partition => '2022-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.block_data_decode';



SELECT partman.create_parent(p_parent_table => 'public.aa_block_info',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'monthly',
                             p_premake => 5,
                             p_start_partition => '2023-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_block_info';





SELECT partman.create_parent(p_parent_table => 'public.transaction_decode',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'daily',
                             p_premake => 30,
                             p_start_partition => '2022-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.transaction_decode';





SELECT partman.create_parent(p_parent_table => 'public.transaction_receipt_decode',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'daily',
                             p_premake => 30,
                             p_start_partition => '2022-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.transaction_receipt_decode';







SELECT partman.create_parent(
               p_parent_table=>'public.aa_transaction_info',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'weekly',
               p_premake=>5,
               p_start_partition=>'2023-01-01'
       );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_transaction_info';




SELECT partman.create_parent(
               p_parent_table=>'public.aa_user_ops_calldata',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'weekly',
               p_premake=>5,
               p_start_partition=>'2023-01-01'
       );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_user_ops_calldata';



SELECT partman.create_parent(
               p_parent_table=>'public.aa_user_ops_info',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'weekly',
               p_premake=>5,
               p_start_partition=>'2023-01-01'
       );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_user_ops_info';


INSERT INTO public.block_sync (block_num, scanned, create_time, update_time) VALUES (14502, false, '2024-01-10 09:41:27.088000 +00:00', '2024-01-10 09:41:31.477000 +00:00');
